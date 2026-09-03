package models

import "time"

type Building struct {
	ID string `gorm:"primaryKey" json:"id"`

	Name string `gorm:"size:100;not null;uniqueIndex" json:"name"`

	Floors []Floor `gorm:"foreignKey:BuildingID" json:"floors,omitempty"`
	Users  []User  `gorm:"foreignKey:BuildingID" json:"users,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

