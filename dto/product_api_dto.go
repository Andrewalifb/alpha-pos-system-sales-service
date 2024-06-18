package dto

type PosProduct struct {
	ProductID          string   `json:"product_id"`
	ProductBarcodeId   string   `json:"product_barcode_id"`
	ProductName        string   `json:"product_name"`
	Price              float64  `json:"price"`
	CostPrice          float64  `json:"cost_price"`
	CategoryID         string   `json:"category_id"`
	SubCategoryID      string   `json:"sub_category_id"`
	StockQuantity      int32    `json:"stock_quantity"`
	ReorderLevel       int32    `json:"reorder_level"`
	SupplierID         string   `json:"supplier_id"`
	ProductDescription string   `json:"product_description"`
	Active             bool     `json:"active"`
	BranchID           string   `json:"branch_id"`
	CompanyID          string   `json:"company_id"`
	CreatedAt          JSONTime `json:"created_at"`
	CreatedBy          string   `json:"created_by"`
	UpdatedAt          JSONTime `json:"updated_at"`
	UpdatedBy          string   `json:"updated_by"`
}
