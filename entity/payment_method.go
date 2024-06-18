package entity

import (
	"time"

	"github.com/google/uuid"
)

type PosPaymentMethod struct {
	PaymentMethodID uuid.UUID `gorm:"type:uuid;primary_key" json:"payment_method_id"`
	MethodName      string    `gorm:"type:varchar(255);not null" json:"method_name"`
	CompanyID       uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt       time.Time `gorm:"type:timestamp" json:"created_at"`
	CreatedBy       uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedAt       time.Time `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy       uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}
