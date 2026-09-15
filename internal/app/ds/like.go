package ds

type Like struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null"`
	PigmentID uint `gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
	Pigment   Pigment `gorm:"foreignKey:PigmentID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}
