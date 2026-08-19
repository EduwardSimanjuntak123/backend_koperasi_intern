package models

type CartItem struct {
    ID        uint    `gorm:"primaryKey" json:"id"`
    
    // Relasi ke Cart (Keranjang Induk)
    CartID    uint    `gorm:"not null" json:"cart_id"`
    Cart      Cart    `gorm:"foreignKey:CartID" json:"cart"`
    
    // Relasi ke Product
    ProductID uint    `gorm:"not null" json:"product_id"`
    Product   Product `gorm:"foreignKey:ProductID" json:"product"`
    
    // Kuantitas barang yang dimasukkan ke keranjang
    Quantity  int     `gorm:"not null;default:1" json:"quantity"`
}
