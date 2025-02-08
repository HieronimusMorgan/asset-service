package assets

// AssetMaintenanceResponse is a struct that represents the response of asset maintenance
type AssetMaintenanceResponse struct {
	ID                 uint    `json:"id"`
	AssetID            int     `json:"asset_id"`
	MaintenanceDate    string  `json:"maintenance_date"`
	MaintenanceDetails *string `json:"maintenance_details,omitempty"`
	MaintenanceCost    float64 `json:"maintenance_cost"`
}
