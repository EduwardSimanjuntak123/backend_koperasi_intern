package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Username  string    `gorm:"not null" json:"username"`
	Email     string    `gorm:"not null" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	No_hp     string    `gorm:"not null" json:"no_hp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Stores    []Store   `gorm:"foreignKey:UserID" json:"stores,omitempty"`
	Favorites []Product `gorm:"many2many:user_favorites;" json:"favorites"`
}

func (User) TableName() string {
	return "user"
}
