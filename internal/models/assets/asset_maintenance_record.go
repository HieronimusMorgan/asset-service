package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceRecord struct {
	MaintenanceRecordID uint            `gorm:"primaryKey" json:"maintenance_record_id"`
	UserClientID        string          `gorm:"type:varchar(50);not null" json:"user_client_id,omitempty"`
	AssetID             int             `gorm:"not null" json:"asset_id"`
	MaintenanceID       int             `gorm:"not null" json:"maintenance_id"`
	MaintenanceTypeID   int             `gorm:"not null" json:"maintenance_type_id"`
	MaintenanceDate     *time.Time      `gorm:"type:date;not null" json:"maintenance_date"`
	MaintenanceDetails  *string         `gorm:"type:text" json:"maintenance_details,omitempty"`
	MaintenanceCost     float64         `gorm:"type:decimal(15,2)" json:"maintenance_cost"`
	PerformedBy         *string         `gorm:"type:text" json:"performed_by,omitempty"`
	IntervalDays        *int            `gorm:"type:int" json:"interval_days,omitempty"`
	NextDueDate         *time.Time      `gorm:"type:date" json:"next_due_date,omitempty"`
	CreatedAt           *time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy           *string         `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedAt           *time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy           *string         `gorm:"type:varchar(255)" json:"updated_by"`
	DeletedAt           *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy           *string         `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
