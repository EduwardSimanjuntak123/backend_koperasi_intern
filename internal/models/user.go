package models

import "time"

type User struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	Nama          string   `gorm:"not null" json:"nama"`
	Username	  varchar  `gorm:"not null" json:"username"`
	Email		  varchar  `gorm:"not null" json:"email"`
	Password	  varchar  `gorm:"not null" json:"password"`
	No_hp		  varchar  `gorm:"not null" json:"no_hp"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}0