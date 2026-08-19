package models

import "time"

type ProductBadge string

const (
	BadgeNew        ProductBadge = "NEW"
	BadgeBestSeller ProductBadge = "BEST_SELLER"
)

type InventoryMovement string

const (
	FastMoving InventoryMovement = "FAST_MOVING"
	SlowMoving InventoryMovement = "SLOW_MOVING"
)

type Product struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Basic
	Name    string  `gorm:"not null" json:"name"`
	Slug    string  `gorm:"unique;not null" json:"slug"`
	Barcode string  `gorm:"size:100;uniqueIndex" json:"barcode"`
	SKU     string  `gorm:"size:100;uniqueIndex" json:"sku"`
	Image   *string `json:"image"`

	// Price
	Price          float64  `gorm:"not null" json:"price"`           // Selling Price
	PurchasePrice  float64  `gorm:"default:0" json:"purchase_price"` // Harga beli
	OriginalPrice  *float64 `json:"original_price"`
	PromotionPrice *float64 `json:"promotion_price"`

	// Stock
	Stock    int  `gorm:"not null;default:0" json:"stock"`
	MinStock *int `json:"min_stock"`
	MaxStock *int `json:"max_stock"`

	// Inventory
	InventoryMovement InventoryMovement `gorm:"default:'FAST_MOVING'" json:"inventory_movement"`
	ExpiredDate       *time.Time        `json:"expired_date"`

	// Dimension
	Volume     string   `json:"volume"`
	WeightGram *int     `json:"weight_gram"`
	LengthCm   *float64 `json:"length_cm"`
	WidthCm    *float64 `json:"width_cm"`
	HeightCm   *float64 `json:"height_cm"`

	Badge *ProductBadge `json:"badge"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relation
	BrandID *uint  `json:"brand_id"`
	Brand   *Brand `gorm:"foreignKey:BrandID" json:"brand"`

	UnitID     *uint            `json:"unit_id"`
	Unit       *Unit            `gorm:"foreignKey:UnitID" json:"unit"`
	CategoryID *uint            `json:"category_id"`
	Category   *CategoryProduct `gorm:"foreignKey:CategoryID" json:"category"`

	Favorites []Product `gorm:"many2many:user_favorites;"`

	StoreID uint  `gorm:"not null" json:"store_id"`
	Store   Store `gorm:"foreignKey:StoreID" json:"store"`
}
