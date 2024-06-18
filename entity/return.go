package entity

import (
	"time"

	"github.com/google/uuid"
)

type PosReturn struct {
	ReturnID   uuid.UUID `gorm:"type:uuid;primary_key" json:"return_id"`
	ReceiptID  string    `gorm:"not null" json:"receipt_id"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity   int       `gorm:"type:int;not null" json:"quantity"`
	Price      float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Amount     float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	ReturnDate time.Time `gorm:"type:timestamp;not null" json:"return_date"`
	Reason     string    `gorm:"type:text" json:"reason"`
	StoreID    uuid.UUID `gorm:"type:uuid" json:"store_id"`
	BranchID   uuid.UUID `gorm:"type:uuid;not null" json:"branch_id"`
	CompanyID  uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt  time.Time `gorm:"type:timestamp" json:"created_at"`
	CreatedBy  uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedAt  time.Time `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy  uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}
