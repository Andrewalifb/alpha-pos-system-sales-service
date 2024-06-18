package dto

import (
	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
)

type GetNextReceiptIDApiRequest struct {
	StoreID string `json:"store_id"`
	*pb.JWTPayload
}

type GetNextReceiptIDApiResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ReceiptID int `json:"receipt_id"`
	} `json:"data"`
}
