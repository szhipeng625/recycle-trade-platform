package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Role      string    `gorm:"size:20;not null;default:user" json:"role"` // admin / user
	Name      string    `gorm:"size:50" json:"name"`
	Phone     string    `gorm:"size:20" json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type Goods struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Category  string    `gorm:"size:50" json:"category"`
	Unit      string    `gorm:"size:20;not null;default:公斤" json:"unit"` // 公斤/市斤/吨
	CreatedAt time.Time `json:"created_at"`
}

type Pricing struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GoodsID   uint      `gorm:"index;not null" json:"goods_id"`
	Goods     Goods     `gorm:"foreignKey:GoodsID" json:"goods,omitempty"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"` // 元/公斤
	Date      string    `gorm:"size:10;index;not null" json:"date"`       // YYYY-MM-DD
	CreatedAt time.Time `json:"created_at"`
}

type Seller struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Phone     string    `gorm:"size:20" json:"phone"`
	IdCard    string    `gorm:"size:30" json:"id_card"`
	Address   string    `gorm:"size:255" json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SellerID       uint      `gorm:"index;not null" json:"seller_id"`
	Seller         Seller    `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	GoodsID        uint      `gorm:"index;not null" json:"goods_id"`
	Goods          Goods     `gorm:"foreignKey:GoodsID" json:"goods,omitempty"`
	Weight         float64   `gorm:"type:decimal(12,3);not null" json:"weight"`
	Unit           string    `gorm:"size:20;not null" json:"unit"` // 吨/公斤/市斤
	WeightInKg     float64   `gorm:"type:decimal(12,3);not null" json:"weight_in_kg"`
	UnitPrice      float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"` // 元/公斤
	TotalAmount    float64   `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	PricingID      uint      `gorm:"index" json:"pricing_id"`
	OperatorID     uint      `gorm:"index" json:"operator_id"`
	Operator       User      `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
	TransactionTime time.Time `json:"transaction_time"`
	CreatedAt      time.Time `json:"created_at"`
}

type Inventory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GoodsID   uint      `gorm:"uniqueIndex;not null" json:"goods_id"`
	Goods     Goods     `gorm:"foreignKey:GoodsID" json:"goods,omitempty"`
	Quantity  float64   `gorm:"type:decimal(12,3);not null;default:0" json:"quantity"` // 公斤
	UpdatedAt time.Time `json:"updated_at"`
}

type InventoryRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	GoodsID       uint      `gorm:"index;not null" json:"goods_id"`
	Goods         Goods     `gorm:"foreignKey:GoodsID" json:"goods,omitempty"`
	Type          string    `gorm:"size:10;not null" json:"type"` // in / out
	Quantity      float64   `gorm:"type:decimal(12,3);not null" json:"quantity"` // 公斤
	Unit          string    `gorm:"size:20;not null" json:"unit"`
	ReferenceType string    `gorm:"size:30" json:"reference_type"` // transaction / manual
	ReferenceID   uint      `gorm:"index" json:"reference_id"`
	Remark        string    `gorm:"size:255" json:"remark"`
	OperatorID    uint      `gorm:"index" json:"operator_id"`
	Operator      User      `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
