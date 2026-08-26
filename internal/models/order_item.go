package models

type Order_Item struct {
	ID      string `gorm:"primaryKey" json:"id"`
	UserID  string `gorm:"not null" json:"user_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user"`
	ChartID string `gorm:"not null" json:"cart_id"`
	Chart   Cart   `gorm:"foreignKey:ChartID" json:"cart"`
}
