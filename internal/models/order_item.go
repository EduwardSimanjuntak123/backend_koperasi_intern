package models

type Order_Item struct {
	ID      uint `gorm:"primaryKey" json:"id"`
	UserID  uint `gorm:"not null" json:"user_id"`
	User    User `gorm:"foreignKey:UserID" json:"user"`
	ChartID uint `gorm:"not null" json:"cart_id"`
	Chart   Cart `gorm:"foreignKey:ChartID" json:"cart"`
}
