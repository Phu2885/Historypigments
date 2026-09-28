package ds

import (
	"strings"
	"time"
)

const minioBase = "/minio/historypigments/"

type Pigment struct {
	ID               uint      `gorm:"primaryKey"`
	Name             string    `gorm:"type:varchar(100);not null"`
	ShortDescription string    `gorm:"type:text"`
	Status           string    `gorm:"type:varchar(20);not null;default:'черновик'"` // черновик/опубликован/удалён
	ImageURL         string    `gorm:"type:varchar(255)"`
	VideoURL         string    `gorm:"type:varchar(255)"`
	Color            string    `gorm:"type:varchar(50)"`
	Reflection       float64   `gorm:"type:decimal(5,2)"`
	Absorption       float64   `gorm:"type:decimal(5,2)"`
	CreatedAt        time.Time `gorm:"not null"`
	FormedAt         *time.Time
	CreatorID        uint
	Creator          User `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (p Pigment) GetImage() string {
	img := strings.TrimSpace(p.ImageURL)
	if img == "" {
		return "/assets/images/placeholder.svg"
	}
	if strings.HasPrefix(img, "/assets/") || strings.HasPrefix(img, "/minio/") {
		return img
	}
	if strings.HasPrefix(img, "http://localhost:9000/historypigments/") {
		return strings.Replace(img, "http://localhost:9000/historypigments/", minioBase, 1)
	}
	if strings.HasPrefix(img, "http") {
		return img
	}
	return minioBase + img
}

func (p Pigment) GetVideo() string {
	vid := strings.TrimSpace(p.VideoURL)
	if vid == "" {
		return ""
	}
	if strings.HasPrefix(vid, "/assets/") || strings.HasPrefix(vid, "/minio/") {
		return vid
	}
	if strings.HasPrefix(vid, "http://localhost:9000/historypigments/") {
		return strings.Replace(vid, "http://localhost:9000/historypigments/", minioBase, 1)
	}
	if strings.HasPrefix(vid, "http") {
		return vid
	}
	return minioBase + vid
}
