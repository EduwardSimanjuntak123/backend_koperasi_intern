package models

import "time"

type Favorite struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID string `gorm:"not null;uniqueIndex:idx_user_product" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`

	ProductID string `gorm:"not null;uniqueIndex:idx_user_product" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (Favorite) TableName() string {
	return "favorites"
}
