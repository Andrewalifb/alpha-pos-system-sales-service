package dto

type PosPromotion struct {
	PromotionId  string   `json:"promotion_id"`
	ProductId    string   `json:"product_id"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Active       bool     `json:"active"`
	DiscountRate float64  `json:"discount_rate"`
	StoreId      string   `json:"store_id"`
	BranchId     string   `json:"branch_id"`
	CompanyId    string   `json:"company_id"`
	CreatedAt    JSONTime `json:"created_at"`
	CreatedBy    string   `json:"created_by"`
	UpdatedAt    JSONTime `json:"updated_at"`
	UpdatedBy    string   `json:"updated_by"`
}

type ReadPosPromotionResponse struct {
	PosPromotion *PosPromotion `json:"pos_promotion"`
}

type ReadPromotionApiSuccessResponse struct {
	Status  bool                      `json:"status"`
	Message string                    `json:"message"`
	Data    *ReadPosPromotionResponse `json:"data"`
}
