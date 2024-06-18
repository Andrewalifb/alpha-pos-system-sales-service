package dto

import (
	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
)

type ReadProductApiRequest struct {
	ProductID string `json:"product_id"`
	*pb.JWTPayload
}

type ReadProductApiResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		PosProduct `json:"pos_product"`
	} `json:"data"`
}
