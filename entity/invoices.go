package entity

import (
	"time"

	"github.com/google/uuid"
)

type PosInvoice struct {
	InvoiceID uuid.UUID `gorm:"type:uuid;primary_key" json:"invoice_id"`
	ReceiptID string    `gorm:"not null" json:"receipt_id"`
	Date      time.Time `gorm:"type:timestamp;not null" json:"date"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Discounts float64   `gorm:"type:decimal(10,2)" json:"discounts"`
	Taxes     float64   `gorm:"type:decimal(10,2)" json:"taxes"`
	BranchID  uuid.UUID `gorm:"type:uuid;not null" json:"branch_id"`
	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}
