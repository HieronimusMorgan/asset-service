package in

import "time"

type AssetMaintenanceRequest struct {
	AssetID            int       `gorm:"not null" json:"asset_id"`
	MaintenanceDate    time.Time `gorm:"type:date;not null" json:"maintenance_date"`
	MaintenanceDetails *string   `gorm:"type:text" json:"maintenance_details,omitempty"`
	MaintenanceCost    float64   `gorm:"type:decimal(15,2)" json:"maintenance_cost"`
}
