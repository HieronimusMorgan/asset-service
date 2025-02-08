package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceType struct {
	ID          int64          `json:"id" gorm:"column:type_id;primaryKey"`
	TypeName    string         `json:"type_name" gorm:"unique"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   string         `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   string         `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy   *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
