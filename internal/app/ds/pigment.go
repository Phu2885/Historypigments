package ds

import (
	"strings"
	"time"
)

const minioBase = "http://localhost:9000/historypigments/"

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
	if strings.TrimSpace(p.ImageURL) == "" {
		return "/assets/images/placeholder.jpg"
	}
	if strings.HasPrefix(p.ImageURL, "http") {
		return p.ImageURL
	}
	return minioBase + p.ImageURL
}

func (p Pigment) GetVideo() string {
	if strings.TrimSpace(p.VideoURL) == "" {
		return ""
	}
	if strings.HasPrefix(p.VideoURL, "http") {
		return p.VideoURL
	}
	return minioBase + p.VideoURL
}
