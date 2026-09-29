package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"rip/internal/app/ds"
)

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

// 1. Получение списка опубликованных (ORM)
func (r *Repository) GetAllPublished() ([]ds.Pigment, error) {
	var pigments []ds.Pigment
	err := r.db.Where("status = ?", "опубликован").Find(&pigments).Error
	return pigments, err
}

// 2. Поиск по поглощению (ORM)
func (r *Repository) FilterByAbsorption(val float64) ([]ds.Pigment, error) {
	var pigments []ds.Pigment
	err := r.db.Where("status = ? AND absorption <= ?", "опубликован", val).Find(&pigments).Error
	return pigments, err
}

// 3. Получение черновика пользователя (ORM)
func (r *Repository) GetUserDraft(userID uint) (*ds.Pigment, error) {
	var draft ds.Pigment
	err := r.db.Where("creator_id = ? AND status = ?", userID, "черновик").First(&draft).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &draft, err
}

// 4. Создание новой карточки-черновика (ORM)
func (r *Repository) CreateDraft(p *ds.Pigment) error {
	p.Status = "черновик"
	p.CreatedAt = time.Now()
	return r.db.Create(p).Error
}

// 5. Публикация карточки (ORM)
func (r *Repository) PublishPigment(id uint, desc string, refl, abs float64) error {
	now := time.Now()
	return r.db.Model(&ds.Pigment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":            "опубликован",
		"short_description": desc,
		"reflection":        refl,
		"absorption":        abs,
		"formed_at":         &now,
	}).Error
}

// 6. Логическое удаление через SQL запрос UPDATE (Без ORM)
func (r *Repository) LogicalDelete(id uint) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := "UPDATE pigments SET status = $1 WHERE id = $2"
	_, err = sqlDB.Exec(query, "удалён", id)
	return err
}

// 7. Получение количества лайков (динамический подсчет по коллекции)
func (r *Repository) GetLikesCount(pigmentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&ds.Like{}).Where("pigment_id = ?", pigmentID).Count(&count).Error
	return count, err
}

// ==========================================
// MЕТОДЫ ЛАБОРАТОРНОЙ РАБОТЫ №3 (REST API)
// ==========================================

// GET /api/pigments — получение списка всех пигментов
func (r *Repository) GetAllPigmentsAPI() ([]ds.Pigment, error) {
	var pigments []ds.Pigment
	err := r.db.Preload("Creator").Find(&pigments).Error
	return pigments, err
}

// GET /api/pigments/:id — получение пигмента по ID с авторам (Preload)
func (r *Repository) GetPigmentByID(id uint) (ds.Pigment, error) {
	var pigment ds.Pigment
	err := r.db.Preload("Creator").Where("id = ?", id).First(&pigment).Error
	return pigment, err
}

// POST /api/pigments — добавление пигмента
func (r *Repository) AddPigment(p *ds.Pigment) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.CreatorID == 0 {
		p.CreatorID = 1 // default test user
	}
	return r.db.Create(p).Error
}

// PUT /api/pigments/:id — обновление пигмента
func (r *Repository) ModifyPigment(id uint, p *ds.Pigment) error {
	updateData := map[string]interface{}{
		"name":              p.Name,
		"short_description": p.ShortDescription,
		"reflection":        p.Reflection,
		"absorption":        p.Absorption,
	}
	if p.Status != "" {
		updateData["status"] = p.Status
	}
	if p.ImageURL != "" {
		updateData["image_url"] = p.ImageURL
	}
	if p.VideoURL != "" {
		updateData["video_url"] = p.VideoURL
	}

	return r.db.Model(&ds.Pigment{}).Where("id = ?", id).Updates(updateData).Error
}

// DELETE /api/pigments/:id — жесткое удаление (Hard Delete для ЛР3 API)
func (r *Repository) HardDeletePigment(id uint) error {
	return r.db.Delete(&ds.Pigment{}, id).Error
}

// GET /api/users/:id/pigments — получение пользователя и созданных им пигментов
func (r *Repository) GetUserPigments(userID uint) (*ds.UserPigments, error) {
	var userStocks ds.UserPigments

	err := r.db.Where("id = ?", userID).First(&userStocks.User).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Where("creator_id = ?", userID).Find(&userStocks.Pigments).Error
	if err != nil {
		return nil, err
	}

	return &userStocks, nil
}

// MinIO Upload Method — Сохранение изображения пигмента в MinIO SDK
func (r *Repository) AddOrReplacePigmentImage(pigmentID uint, header *multipart.FileHeader) error {
	if r.minio == nil {
		return fmt.Errorf("minio client not initialized")
	}

	filename := fmt.Sprintf("pigment_%d.png", pigmentID)

	file, err := header.Open()
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	contentType := http.DetectContentType(buffer)

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("ошибка перемещения по файловому потоку: %w", err)
	}

	_, err = r.minio.PutObject(
		context.Background(),
		r.minioBucketName,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf("не удалось добавить объект в хранилище minio: %w", err)
	}

	imagePath := "/minio/" + r.minioBucketName + "/" + filename
	err = r.db.Model(&ds.Pigment{}).Where("id = ?", pigmentID).UpdateColumn("image_url", imagePath).Error
	if err != nil {
		r.minio.RemoveObject(context.Background(), r.minioBucketName, filename, minio.RemoveObjectOptions{})
		return fmt.Errorf("ошибка сохранения пути к изображению в БД: %w", err)
	}

	return nil
}
