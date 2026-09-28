package ds

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Login    string `gorm:"type:varchar(50);unique;not null" json:"login"`
	Password string `gorm:"type:varchar(100);not null" json:"-"`
}

type UserPigments struct {
	User     User      `json:"user"`
	Pigments []Pigment `json:"pigments"`
}
