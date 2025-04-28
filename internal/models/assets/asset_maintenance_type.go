package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceType struct {
	MaintenanceTypeID   uint            `json:"maintenance_type_id" gorm:"column:maintenance_type_id;primaryKey"`
	UserClientID        string          `gorm:"type:varchar(50);not null" json:"user_client_id,omitempty"`
	MaintenanceTypeName string          `json:"maintenance_type_name" gorm:"unique"`
	Description         string          `json:"description"`
	CreatedAt           *time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy           *string         `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedAt           *time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy           *string         `gorm:"type:varchar(255)" json:"updated_by"`
	DeletedAt           *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy           *string         `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
