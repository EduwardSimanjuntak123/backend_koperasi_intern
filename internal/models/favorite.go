package models

type Favorite struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Products []Product `gorm:"many2many:product_categories;" json:"-"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
}
