package in

type AssetRequest struct {
	Name              string  `json:"name" binding:"required"`
	Description       string  `json:"description"`
	CategoryID        int     `json:"category_id" binding:"required"`
	StatusID          int     `json:"status_id" binding:"required"`
	PurchaseDate      string  `json:"purchase_date" binding:"required"`
	ExpiryDate        string  `json:"expiry_date"`
	Value             float64 `json:"value" binding:"required"`
	MaintenanceDate   string  `json:"maintenance_date"`
	MaintenanceCost   float64 `json:"maintenance_cost"`
	MaintenanceDetail string  `json:"maintenance_detail"`
}
