package models

type Cart struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    UserID    string    `gorm:"not null" json:"user_id"`
    User      User      `gorm:"foreignKey:UserID" json:"user"`
    
    // Relasi One-to-Many: Satu keranjang bisa memiliki banyak item
    CartItems []CartItem `gorm:"foreignKey:CartID" json:"cart_items"`
}
