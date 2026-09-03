package models

import "time"

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	NoHP      string    `gorm:"uniqueIndex;not null" json:"no_hp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	RoleID    uint      `gorm:"not null" json:"role_id"`
	Role      Roles     `gorm:"foreignKey:RoleID" json:"role"`

	// Relasi Gedung & Lantai (Opsional / Lokasi Default Pengantaran)
	BuildingID *string   `gorm:"size:50;index" json:"building_id"`
	Building   *Building `gorm:"foreignKey:BuildingID" json:"building,omitempty"`
	FloorID    *string   `gorm:"size:50;index" json:"floor_id"`
	Floor      *Floor    `gorm:"foreignKey:FloorID" json:"floor,omitempty"`

	Stores    []Store   `gorm:"foreignKey:UserID" json:"stores,omitempty"`
	Favorites []Product `gorm:"many2many:user_favorites;" json:"favorites"`
}

func (User) TableName() string {
	return "users"
}

