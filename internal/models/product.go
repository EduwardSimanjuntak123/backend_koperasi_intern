package models

import "time"

type ProductBadge string

const (
	BadgeNew        ProductBadge = "NEW"
	BadgeBestSeller ProductBadge = "BEST_SELLER"
)

type Product struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	Name          string   `gorm:"not null" json:"name"`
	Slug          string   `gorm:"unique;not null" json:"slug"`
	Image         *string  `json:"image"`
	Price         float64  `gorm:"not null" json:"price"`
	OriginalPrice *float64 `json:"original_price"`
	Stock         int      `gorm:"not null;default:0" json:"stock"`

	Volume     string   `json:"volume"`
	WeightGram *int     `json:"weight_gram"`
	LengthCm   *float64 `json:"length_cm"`
	WidthCm    *float64 `json:"width_cm"`
	HeightCm   *float64 `json:"height_cm"`

	Badge *ProductBadge `json:"badge"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Categories []Category `gorm:"many2many:product_categories;" json:"categories"`
}
