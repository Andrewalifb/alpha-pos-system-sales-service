package dto

import pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"

// PosInventoryHistory represents the inventory history information
type PosInventoryHistory struct {
	InventoryId string   `json:"inventory_id"`
	ProductId   string   `json:"product_id"`
	StoreId     string   `json:"store_id"`
	Date        JSONTime `json:"date"`
	Quantity    int32    `json:"quantity"`
	BranchId    string   `json:"branch_id"`
	CompanyId   string   `json:"company_id"`
	CreatedAt   JSONTime `json:"created_at"`
	CreatedBy   string   `json:"created_by"`
	UpdatedAt   JSONTime `json:"updated_at"`
	UpdatedBy   string   `json:"updated_by"`
}

// CreatePosInventoryHistoryRequest is the request structure for creating inventory history
type CreatePosInventoryHistoryRequest struct {
	PosInventoryHistory *PosInventoryHistory `json:"pos_inventory_history"`
	JwtPayload          *pb.JWTPayload       `json:"jwt_payload"`
	JwtToken            string               `json:"jwt_token"`
}

// CreatePosInventoryHistoryResponse is the response structure after creating inventory history
type CreatePosInventoryHistoryResponse struct {
	PosInventoryHistory *PosInventoryHistory `json:"pos_inventory_history"`
}

// SuccessResponse is the standard structure for a successful API response
type CreateInventorySuccessResponse struct {
	Status  bool                               `json:"status"`
	Message string                             `json:"message"`
	Data    *CreatePosInventoryHistoryResponse `json:"data"`
}
