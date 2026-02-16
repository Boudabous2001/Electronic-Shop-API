package models

import (
	"net/url"
	"time"
)

// ========================================
// 🏪 SHOP - Représente une boutique
// ========================================
type Shop struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	Active         bool      `gorm:"default:true" json:"active"`
	WhatsAppNumber string    `gorm:"not null" json:"whatsapp_number"`
	CreatedAt      time.Time `json:"created_at"`
}

// ========================================
// 👤 USER - Utilisateurs du système
// ========================================
type Role string

const (
	RoleSuperAdmin Role = "SuperAdmin"
	RoleAdmin      Role = "Admin"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // json:"-" = jamais exposé
	Role      Role      `gorm:"not null" json:"role"`
	ShopID    uint      `gorm:"not null" json:"shop_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ========================================
// 📦 PRODUCT - Produits du magasin
// ========================================
type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	PurchasePrice float64   `gorm:"not null" json:"purchase_price,omitempty"`
	SellingPrice  float64   `gorm:"not null" json:"selling_price"`
	Stock         int       `gorm:"default:0" json:"stock"`
	ImageURL      string    `json:"image_url"`
	ShopID        uint      `gorm:"not null" json:"shop_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// ProductPublic - Version publique sans PurchasePrice
type ProductPublic struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	SellingPrice float64 `json:"selling_price"`
	Stock        int     `json:"stock"`
	ImageURL     string  `json:"image_url"`
	InStock      bool    `json:"in_stock"`
	WhatsAppLink string  `json:"whatsapp_link"`
}

// ToPublic convertit un Product en ProductPublic
func (p *Product) ToPublic(whatsappNumber string) ProductPublic {
	return ProductPublic{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		Category:     p.Category,
		SellingPrice: p.SellingPrice,
		Stock:        p.Stock,
		ImageURL:     p.ImageURL,
		InStock:      p.Stock > 0,
		WhatsAppLink: GenerateWhatsAppLink(whatsappNumber, p.Name),
	}
}

// ========================================
// 💰 TRANSACTION
// ========================================
type TransactionType string

const (
	TypeSale       TransactionType = "Sale"
	TypeExpense    TransactionType = "Expense"
	TypeWithdrawal TransactionType = "Withdrawal"
)

type Transaction struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Type      TransactionType `gorm:"not null" json:"type"`
	ProductID *uint           `json:"product_id,omitempty"`
	Quantity  int             `json:"quantity"`
	Amount    float64         `gorm:"not null" json:"amount"`
	ShopID    uint            `gorm:"not null" json:"shop_id"`
	CreatedAt time.Time       `json:"created_at"`
	Product   *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// ========================================
// 📱 WHATSAPP LINK GENERATOR
// ========================================
func GenerateWhatsAppLink(whatsappNumber string, productName string) string {
	message := "Bonjour je veux plus d'information sur " + productName
	encodedMessage := url.QueryEscape(message)
	return "https://wa.me/" + whatsappNumber + "?text=" + encodedMessage
}