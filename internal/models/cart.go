package models

type Cart struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`
}
