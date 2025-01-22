package out

import "time"

type AssetResponse struct {
	ID           uint      `json:"asset_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	CategoryName string    `json:"category_name,omitempty"`
	StatusName   string    `json:"status_name,omitempty"`
	PurchaseDate time.Time `json:"purchase_date,omitempty"`
	Value        float64   `json:"value,omitempty"`
}
