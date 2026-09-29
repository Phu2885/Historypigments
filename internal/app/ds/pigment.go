package ds

import (
	"strings"
	"time"
)

const minioBase = "/minio/historypigments/"

type Pigment struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Name             string     `gorm:"type:varchar(100);not null" json:"name"`
	ShortDescription string     `gorm:"type:text" json:"short_description"`
	Status           string     `gorm:"type:varchar(20);not null;default:'черновик'" json:"status"` // черновик/опубликован/удалён
	ImageURL         string     `gorm:"type:varchar(255)" json:"image_url"`
	VideoURL         string     `gorm:"type:varchar(255)" json:"video_url"`
	Reflection       float64    `gorm:"type:decimal(5,2)" json:"reflection"`
	Absorption       float64    `gorm:"type:decimal(5,2)" json:"absorption"`
	CreatedAt        time.Time  `gorm:"not null" json:"created_at"`
	FormedAt         *time.Time `json:"formed_at,omitempty"`
	CreatorID        uint       `json:"creator_id"`
	Creator          User       `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	IsCreator        int        `gorm:"-" json:"is_creator"` // 0/1 признак совпадения создателя с текущим пользователем
	LikesCount       int64      `gorm:"-" json:"likes_count,omitempty"`
}

type FullPigmentSerializer struct {
	Pigment
	Creator User `json:"creator"`
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
