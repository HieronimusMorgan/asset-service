package assets

type AssetMaintenanceRequest struct {
	AssetID            int     `gorm:"not null" json:"asset_id"`
	TypeID             int     `gorm:"not null" json:"type_id"`
	MaintenanceDate    string  `gorm:"type:date;not null" json:"maintenance_date"`
	MaintenanceDetails *string `gorm:"type:text" json:"maintenance_details,omitempty"`
	MaintenanceCost    float64 `gorm:"type:decimal(15,2)" json:"maintenance_cost"`
	NextDueDate        string  `gorm:"type:date" json:"next_due_date,omitempty"`
}
