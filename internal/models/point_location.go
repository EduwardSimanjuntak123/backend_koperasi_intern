package models

type PointLocation struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	StoreID       uint   `gorm:"not null" json:"store_id"`
	Store         Store  `gorm:"foreignKey:StoreID" json:"store"`
	Name_Location string `gorm:"not null" json:"name"`
	address       string `gorm:"not null" json:"address"`
	location      string `gorm:"not null" json:"location"`
	shipping_cost int    `gorm:"not null" json:"ongkos_kirim"`
	latitude      string `gorm:"not null" json:"latitude"`
	longitude     string `gorm:"not null" json:"longitude"`
}
