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
	ID                 uint                    `json:"id"`
	UserClientID       string                  `json:"user_client_id,omitempty"`
	AssetID            int                     `json:"asset_id"`
	Type               MaintenanceTypeResponse `json:"type"`
	MaintenanceDate    *DateOnly               `json:"maintenance_date"`
	MaintenanceDetails *string                 `json:"maintenance_details,omitempty"`
	MaintenanceCost    float64                 `json:"maintenance_cost"`
	PerformedBy        *string                 `json:"performed_by,omitempty"`
	IntervalDays       *int                    `json:"interval_days,omitempty"`
	NextDueDate        *DateOnly               `json:"next_due_date,omitempty"`
}

type AssetMaintenancesWithMaintenanceRecordResponse struct {
	ID                 uint                              `json:"id"`
	UserClientID       string                            `json:"user_client_id,omitempty"`
	AssetID            int                               `json:"asset_id"`
	MaintenanceDate    *DateOnly                         `json:"maintenance_date"`
	MaintenanceDetails *string                           `json:"maintenance_details,omitempty"`
	MaintenanceCost    float64                           `json:"maintenance_cost"`
	PerformedBy        *string                           `json:"performed_by,omitempty"`
	IntervalDays       *int                              `json:"interval_days,omitempty"`
	NextDueDate        *DateOnly                         `json:"next_due_date,omitempty"`
	MaintenanceRecord  *[]AssetMaintenanceRecordResponse `json:"maintenance_record,omitempty"`
}

type MaintenanceTypeResponse struct {
	MaintenanceTypeID   int    `json:"maintenance_type_id"`
	MaintenanceTypeName string `json:"maintenance_type_name"`
}

type AssetMaintenanceRecordResponse struct {
	MaintenanceRecordID uint       `json:"maintenance_record_id"`
	MaintenanceTypeName string     `json:"maintenance_type_name"`
	MaintenanceDetails  *string    `json:"maintenance_details,omitempty"`
	MaintenanceDate     *time.Time `json:"maintenance_date"`
	MaintenanceCost     float64    `json:"maintenance_cost"`
	PerformedBy         *string    `json:"performed_by,omitempty"`
	IntervalDays        *int       `json:"interval_days,omitempty"`
	NextDueDate         *time.Time `json:"next_due_date,omitempty"`
}
