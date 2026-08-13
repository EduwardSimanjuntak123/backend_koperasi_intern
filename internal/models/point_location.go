package models

type PointLocation struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	StoreID       uint   `gorm:"not null" json:"store_id"`
	Store         Store  `gorm:"foreignKey:StoreID" json:"store"`
	Name_Location string `gorm:"not null" json:"name"`
	Address       string `gorm:"not null" json:"address"`
	Location      string `gorm:"not null" json:"location"`
	Shipping_Cost int    `gorm:"not null" json:"ongkos_kirim"`
	Latitude      string `gorm:"not null" json:"latitude"`
	Longitude     string `gorm:"not null" json:"longitude"`
}
