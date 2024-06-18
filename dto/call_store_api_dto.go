package dto

import (
	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
)

type PosStoreApiRequest struct {
	StoreID string `json:"store_id"`
	*pb.JWTPayload
}

type PosStoreApiResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		PosStore `json:"pos_store"`
	} `json:"data"`
}

type PosStore struct {
	StoreId   string   `json:"store_id"`
	StoreName string   `json:"store_name"`
	BranchId  string   `json:"branch_id"`
	Location  string   `json:"location"`
	CompanyId string   `json:"company_id"`
	CreatedAt JSONTime `json:"created_at"`
	CreatedBy string   `json:"created_by"`
	UpdatedAt JSONTime `json:"updated_at"`
	UpdatedBy string   `json:"updated_by"`
}
