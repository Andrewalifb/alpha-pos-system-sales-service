package entity

import (
	"time"

	"github.com/google/uuid"
)

type PosCustomer struct {
	CustomerID       uuid.UUID `gorm:"type:uuid;primary_key" json:"customer_id"`
	FirstName        string    `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName         string    `gorm:"type:varchar(255);not null" json:"last_name"`
	Email            string    `gorm:"type:varchar(255)" json:"email"`
	PhoneNumber      string    `gorm:"type:varchar(20)" json:"phone_number"`
	DateOfBirth      time.Time `gorm:"type:date" json:"date_of_birth"`
	RegistrationDate time.Time `gorm:"type:date" json:"registration_date"`
	Address          string    `gorm:"type:varchar(255)" json:"address"`
	City             string    `gorm:"type:varchar(100)" json:"city"`
	Country          string    `gorm:"type:varchar(100)" json:"country"`
	BranchID         uuid.UUID `gorm:"type:uuid;not null" json:"branch_id"`
	CompanyID        uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt        time.Time `gorm:"type:timestamp" json:"created_at"`
	CreatedBy        uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedAt        time.Time `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy        uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}
