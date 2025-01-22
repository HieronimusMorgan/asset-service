package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetStatus struct {
	AssetStatusID uint           `gorm:"primaryKey" json:"asset_status_id"`
	StatusName    string         `gorm:"type:varchar(50);not null" json:"status_name"`
	Description   string         `gorm:"type:text" json:"description"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy     string         `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy     string         `gorm:"type:varchar(50)" json:"updated_by"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy     *string        `gorm:"type:varchar(50)" json:"deleted_by,omitempty"`
}
