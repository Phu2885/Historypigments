package repository

import (
	"gorm.io/gorm"
	"rip/internal/app/ds"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
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
