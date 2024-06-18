package dto

import (
	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
)

type PosUserApiRequest struct {
	UserID string `json:"user_id"`
	*pb.JWTPayload
}

type PosUserApiResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		PosUser `json:"pos_user"`
	} `json:"data"`
}

type PosUser struct {
	UserId       string   `json:"user_id"`
	Username     string   `json:"username"`
	PasswordHash string   `json:"password_hash"`
	RoleId       string   `json:"role_id"`
	CompanyId    string   `json:"company_id"`
	BranchId     string   `json:"branch_id"`
	StoreId      string   `json:"store_id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email"`
	PhoneNumber  string   `json:"phone_number"`
	CreatedAt    JSONTime `json:"created_at"`
	CreatedBy    string   `json:"created_by"`
	UpdatedAt    JSONTime `json:"updated_at"`
	UpdatedBy    string   `json:"updated_by"`
}
