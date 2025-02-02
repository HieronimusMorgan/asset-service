package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceRecord struct {
	MaintenanceRecordID uint           `gorm:"primaryKey" json:"maintenance_record_id"`
	AssetID             uint           `gorm:"not null" json:"asset_id"`
	MaintenanceDate     time.Time      `gorm:"type:date;not null" json:"maintenance_date"`
	MaintenanceType     string         `gorm:"type:varchar(255);not null" json:"maintenance_type"`
	MaintenanceDetails  *string        `gorm:"type:text" json:"maintenance_details,omitempty"`
	MaintenanceCost     *float64       `gorm:"type:decimal(15,2)" json:"maintenance_cost,omitempty"`
	PerformedBy         *string        `gorm:"type:varchar(255)" json:"performed_by,omitempty"`
	NextDueDate         *time.Time     `gorm:"type:date" json:"next_due_date,omitempty"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy           string         `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy           string         `gorm:"type:varchar(255)" json:"updated_by"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy           *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
