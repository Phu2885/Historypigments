package ds

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(100);not null"`
	Login    string `gorm:"type:varchar(50);unique;not null"`
	Password string `gorm:"type:varchar(100);not null"`
}
