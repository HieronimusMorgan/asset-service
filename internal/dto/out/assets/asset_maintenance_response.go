package assets

// AssetMaintenanceResponse is a struct that represents the response of asset maintenance
type AssetMaintenanceResponse struct {
	ID                 uint    `json:"id"`
	AssetID            int     `json:"asset_id"`
	MaintenanceDate    string  `json:"maintenance_date"`
	MaintenanceDetails *string `json:"maintenance_details,omitempty"`
	MaintenanceCost    float64 `json:"maintenance_cost"`
}

type AssetMaintenancesResponse struct {
	ID                 uint                    `gorm:"primaryKey" json:"id"`
	UserClientID       string                  `gorm:"type:varchar(50);not null" json:"user_client_id,omitempty"`
	AssetID            int                     `gorm:"not null" json:"asset_id"`
	Type               MaintenanceTypeResponse `gorm:"not null" json:"type"`
	MaintenanceDate    *DateOnly               `gorm:"type:date;not null" json:"maintenance_date"`
	MaintenanceDetails *string                 `gorm:"type:text" json:"maintenance_details,omitempty"`
	MaintenanceCost    float64                 `gorm:"type:decimal(15,2)" json:"maintenance_cost"`
	PerformedBy        *string                 `gorm:"type:text" json:"performed_by,omitempty"`
	IntervalDays       *int                    `gorm:"type:int" json:"interval_days,omitempty"`
	NextDueDate        *DateOnly               `gorm:"type:date" json:"next_due_date,omitempty"`
}

type MaintenanceTypeResponse struct {
	TypeID   int    `json:"type_id"`
	TypeName string `json:"type_name"`
}
