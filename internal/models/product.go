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

type ProductFilter struct {
	Search            string
	BrandID           *string
	CategoryID        *string
	UnitID            *string
	StoreID           *string
	MinPrice          *float64
	MaxPrice          *float64
	InventoryMovement *InventoryMovement
	Badge             *ProductBadge
	StockStatus       string
	Expired           string
	CreatedFrom       *time.Time
	CreatedTo         *time.Time

	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

type ProductSales struct {
	Product   Product `json:"product"`
	TotalSold int     `json:"total_sold"`
}

type Product struct {
	ID        string `gorm:"primaryKey" json:"id"`
	TotalSold int    `gorm:"-" json:"total_sold,omitempty"`

	// Basic
	Name    string  `gorm:"not null" json:"name"`
	Slug    string  `gorm:"unique;not null" json:"slug"`
	Barcode string  `gorm:"size:100;uniqueIndex" json:"barcode"`
	SKU     string  `gorm:"size:100;uniqueIndex" json:"sku"`
	Image   *string `json:"image"`

	// Price
	Price          float64  `gorm:"not null" json:"price"`           // Selling Price
	PurchasePrice  float64  `gorm:"default:0" json:"purchase_price"` // Harga beli
	PromotionPrice *float64 `json:"promotion_price"`

	// Stock
	Stock int `gorm:"not null;default:0" json:"stock"`

	// Inventory
	InventoryMovement *InventoryMovement `gorm:"default:'FAST_MOVING'" json:"inventory_movement"`
	ExpiredDate       *time.Time         `json:"expired_date"`

	Badge *ProductBadge `json:"badge"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relation
	BrandID *string `json:"brand_id"`
	Brand   *Brand  `gorm:"foreignKey:BrandID" json:"brand"`

	UnitID     *string          `json:"unit_id"`
	Unit       *Unit            `gorm:"foreignKey:UnitID" json:"unit"`
	CategoryID *string          `json:"category_id"`
	Category   *CategoryProduct `gorm:"foreignKey:CategoryID" json:"category"`

	Favorites []Product `gorm:"many2many:user_favorites;"`

	StoreID string `gorm:"not null" json:"store_id"`
	Store   Store  `gorm:"foreignKey:StoreID" json:"store"`
}
