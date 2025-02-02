package assets

import (
	"time"

	"gorm.io/gorm"
)

type AssetMaintenance struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	AssetID            int            `gorm:"not null" json:"asset_id"`
	MaintenanceDate    time.Time      `gorm:"type:date;not null" json:"maintenance_date"`
	MaintenanceDetails *string        `gorm:"type:text" json:"maintenance_details,omitempty"`
	MaintenanceCost    float64        `gorm:"type:decimal(15,2)" json:"maintenance_cost"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy          string         `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy          string         `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy          *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
