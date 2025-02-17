package assets

import "time"

// AssetMaintenanceResponse is a struct that represents the response of asset maintenance
type AssetMaintenanceResponse struct {
	ID                 uint    `json:"id"`
	AssetID            int     `json:"asset_id"`
	MaintenanceDate    string  `json:"maintenance_date"`
	MaintenanceDetails *string `json:"maintenance_details,omitempty"`
	MaintenanceCost    float64 `json:"maintenance_cost"`
}

type AssetMaintenancesResponse struct {
	ID                 uint       `json:"id"`
	TypeName           string     `json:"maintenance_name"`
	MaintenanceDate    time.Time  `json:"maintenance_date"`
	MaintenanceDetails *string    `json:"maintenance_details,omitempty"`
	MaintenanceCost    float64    `json:"maintenance_cost"`
	PerformedBy        *string    `json:"performed_by,omitempty"`
	NextDueDate        *time.Time `json:"next_due_date,omitempty"`
}
