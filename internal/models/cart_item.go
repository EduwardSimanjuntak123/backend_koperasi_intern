package models

type Cart_Item struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`
	CartID uint `gorm:"not null" json:"cart_id"`
	Cart   Cart `gorm:"foreignKey:CartID" json:"cart"`
}
