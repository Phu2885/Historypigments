package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"rip/internal/app/ds"
)

// Singleton current user function
func GetCurrentUserID() uint {
	return 1
}

type Repository struct {
	db              *gorm.DB
	minio           *minio.Client
	minioBucketName string
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db, minioBucketName: "historypigments"}
}

func (r *Repository) SetMinio(minioClient *minio.Client, bucket string) {
	r.minio = minioClient
	if bucket != "" {
		r.minioBucketName = bucket
	}
}

// Helper: Проверка, поставил ли пользователь лайк
func (r *Repository) IsPigmentLikedByUser(pigmentID uint, userID uint) bool {
	var count int64
	r.db.Model(&ds.Like{}).Where("pigment_id = ? AND user_id = ?", pigmentID, userID).Count(&count)
	return count > 0
}

// ==========================================
// MЕТОДЫ ЛАБОРАТОРНОЙ РАБОТЫ №3 (ДОМЕН PIGMENTS & USERS)
// ==========================================

// 1. GET /api/pigments — Список опубликованных пигментов с фильтрацией + признак is_creator (0/1)
func (r *Repository) GetPublishedPigmentsList(absorptionFilter *float64, currentUserID uint) ([]ds.Pigment, error) {
	var pigments []ds.Pigment
	query := r.db.Where("status = ?", "опубликован")
	if absorptionFilter != nil {
		query = query.Where("absorption <= ?", *absorptionFilter)
	}

	err := query.Find(&pigments).Error
	if err != nil {
		return nil, err
	}

	for i := range pigments {
		if pigments[i].CreatorID == currentUserID {
			pigments[i].IsCreator = 1
		} else {
			pigments[i].IsCreator = 0
		}
		likes, _ := r.GetLikesCount(pigments[i].ID)
		pigments[i].Likes = likes
		pigments[i].IsLike = r.IsPigmentLikedByUser(pigments[i].ID, currentUserID)
	}

	return pigments, nil
}

// 2. GET /api/pigments/feed — Лента (возвращает 1 опубликованный пигмент)
func (r *Repository) GetFeedItem(id uint, next bool) (*ds.Pigment, error) {
	var published []ds.Pigment
	err := r.db.Where("status = ?", "опубликован").Order("id asc").Find(&published).Error
	if err != nil || len(published) == 0 {
		return nil, err
	}

	var current ds.Pigment
	if next && id > 0 {
		found := false
		for _, p := range published {
			if p.ID > id {
				current = p
				found = true
				break
			}
		}
		if !found {
			current = published[0]
		}
	} else if id > 0 {
		found := false
		for _, p := range published {
			if p.ID == id {
				current = p
				found = true
				break
			}
		}
		if !found {
			current = published[0]
		}
	} else {
		current = published[0]
	}

	likes, _ := r.GetLikesCount(current.ID)
	current.Likes = likes
	current.IsLike = r.IsPigmentLikedByUser(current.ID, GetCurrentUserID())
	current.IsCreator = 0
	if current.CreatorID == GetCurrentUserID() {
		current.IsCreator = 1
	}

	return &current, nil
}

// 3. GET /api/pigments/draft — Получение черновика текущего пользователя (не более 1 записи)
func (r *Repository) GetUserDraft(userID uint) (*ds.Pigment, error) {
	var draft ds.Pigment
	err := r.db.Where("creator_id = ? AND status = ?", userID, "черновик").First(&draft).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &draft, err
}

// Helper: Генерация латинского имени файла
func sanitizeLatinFilename(original string) string {
	ext := filepath.Ext(original)
	if ext == "" {
		ext = ".bin"
	}
	base := strings.TrimSuffix(original, ext)
	var result strings.Builder
	for _, runeVal := range strings.ToLower(base) {
		if (runeVal >= 'a' && runeVal <= 'z') || (runeVal >= '0' && runeVal <= '9') || runeVal == '_' || runeVal == '-' {
			result.WriteRune(runeVal)
		} else if runeVal == ' ' {
			result.WriteRune('_')
		}
	}
	latin := result.String()
	if latin == "" {
		latin = "file"
	}
	return fmt.Sprintf("%s_%d%s", latin, time.Now().UnixNano()%1000000, ext)
}

// Upload file to MinIO with latin filename
func (r *Repository) UploadFileToMinIO(header *multipart.FileHeader, prefix string) (string, error) {
	if r.minio == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	latinFilename := prefix + "_" + sanitizeLatinFilename(header.Filename)

	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %w", err)
	}

	contentType := http.DetectContentType(buffer)
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("ошибка перемещения по файловому потоку: %w", err)
	}

	_, err = r.minio.PutObject(
		context.Background(),
		r.minioBucketName,
		latinFilename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("не удалось добавить объект в хранилище minio: %w", err)
	}

	return latinFilename, nil
}

// 4. POST /api/pigments — Создание черновика + сохранение файлов в MinIO (латинские имена)
func (r *Repository) CreateDraftPigment(p *ds.Pigment, currentUserID uint, imageHeader *multipart.FileHeader, videoHeader *multipart.FileHeader) error {
	existingDraft, _ := r.GetUserDraft(currentUserID)
	if existingDraft != nil {
		return fmt.Errorf("у вас уже есть активный черновик (ID: %d), опубликуйте или удалите его", existingDraft.ID)
	}

	p.Status = "черновик"
	p.CreatorID = currentUserID
	p.CreatedAt = time.Now()

	if imageHeader != nil {
		imageKey, err := r.UploadFileToMinIO(imageHeader, "img")
		if err == nil {
			p.ImageURL = imageKey
		}
	}

	if videoHeader != nil {
		videoKey, err := r.UploadFileToMinIO(videoHeader, "vid")
		if err == nil {
			p.VideoURL = videoKey
		}
	}

	return r.db.Create(p).Error
}

// 5. PUT /api/pigments/:id/publish — Публикация черновика (смена статуса на 'опубликован')
func (r *Repository) PublishDraftPigment(id uint, currentUserID uint, desc string, refl, abs float64) error {
	var pigment ds.Pigment
	err := r.db.Where("id = ? AND creator_id = ?", id, currentUserID).First(&pigment).Error
	if err != nil {
		return fmt.Errorf("пигмент не найден или у вас нет прав на его изменение")
	}

	if pigment.Status != "черновик" {
		return fmt.Errorf("статус нельзя сменить: пигмент уже опубликован или удален")
	}

	now := time.Now()
	return r.db.Model(&pigment).Updates(map[string]interface{}{
		"status":            "опубликован",
		"short_description": desc,
		"reflection":        refl,
		"absorption":        abs,
		"formed_at":         &now,
	}).Error
}

// 6. DELETE /api/pigments/:id — Soft Delete (только для пигментов текущего пользователя)
func (r *Repository) SoftDeletePigment(id uint, currentUserID uint) error {
	var pigment ds.Pigment
	err := r.db.Where("id = ? AND creator_id = ?", id, currentUserID).First(&pigment).Error
	if err != nil {
		return fmt.Errorf("пигмент не найден или вы не являетесь его создателем")
	}

	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := "UPDATE pigments SET status = $1 WHERE id = $2 AND creator_id = $3"
	_, err = sqlDB.Exec(query, "удалён", id, currentUserID)
	return err
}

// 7. POST /api/pigments/:id/like — Лайк (0 отменяет, 1 ставит)
func (r *Repository) ToggleLike(pigmentID uint, currentUserID uint, value int) (int64, bool, error) {
	if value == 1 {
		var existing ds.Like
		err := r.db.Where("user_id = ? AND pigment_id = ?", currentUserID, pigmentID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			r.db.Create(&ds.Like{UserID: currentUserID, PigmentID: pigmentID})
		}
	} else if value == 0 {
		r.db.Where("user_id = ? AND pigment_id = ?", currentUserID, pigmentID).Delete(&ds.Like{})
	}

	likesCount, err := r.GetLikesCount(pigmentID)
	isLiked := r.IsPigmentLikedByUser(pigmentID, currentUserID)
	return likesCount, isLiked, err
}

// GET /api/pigments/:id — получение пигмента по ID с авторам (Preload)
func (r *Repository) GetPigmentByID(id uint) (ds.Pigment, error) {
	var pigment ds.Pigment
	err := r.db.Preload("Creator").Where("id = ? AND status != ?", id, "удалён").First(&pigment).Error
	if err == nil {
		pigment.Likes, _ = r.GetLikesCount(pigment.ID)
		pigment.IsLike = r.IsPigmentLikedByUser(pigment.ID, GetCurrentUserID())
		if pigment.CreatorID == GetCurrentUserID() {
			pigment.IsCreator = 1
		}
	}
	return pigment, err
}

// GET /api/users/:id/pigments — получение пользователя и созданных им пигментов
func (r *Repository) GetUserPigments(userID uint) (*ds.UserPigments, error) {
	var userPigments ds.UserPigments

	err := r.db.Where("id = ?", userID).First(&userPigments.User).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Where("creator_id = ? AND status != ?", userID, "удалён").Find(&userPigments.Pigments).Error
	if err != nil {
		return nil, err
	}

	for i := range userPigments.Pigments {
		likes, _ := r.GetLikesCount(userPigments.Pigments[i].ID)
		userPigments.Pigments[i].Likes = likes
		userPigments.Pigments[i].IsLike = r.IsPigmentLikedByUser(userPigments.Pigments[i].ID, GetCurrentUserID())
		if userPigments.Pigments[i].CreatorID == GetCurrentUserID() {
			userPigments.Pigments[i].IsCreator = 1
		}
	}

	return &userPigments, nil
}

// GET /api/users — получение всех пользователей
func (r *Repository) GetAllUsers() ([]ds.User, error) {
	var users []ds.User
	err := r.db.Find(&users).Error
	return users, err
}

// Регистрация нового пользователя
func (r *Repository) RegisterUser(u *ds.User) error {
	if u.Role == "" {
		u.Role = "user"
	}
	var existing ds.User
	if err := r.db.Where("login = ?", u.Login).First(&existing).Error; err == nil {
		return fmt.Errorf("пользователь с логином '%s' уже существует", u.Login)
	}
	return r.db.Create(u).Error
}

// Поиск пользователя по логину
func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	var user ds.User
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Получение количества лайков
func (r *Repository) GetLikesCount(pigmentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&ds.Like{}).Where("pigment_id = ?", pigmentID).Count(&count).Error
	return count, err
}
