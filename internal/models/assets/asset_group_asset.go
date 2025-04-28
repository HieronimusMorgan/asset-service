package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetGroupAsset struct {
	AssetID      uint            `gorm:"primaryKey;column:asset_id" json:"asset_id,omitempty"`
	AssetGroupID uint            `gorm:"primaryKey;column:asset_group_id"  json:"asset_group_id,omitempty"`
	UserID       uint            `gorm:"primaryKey;column:user_id"  json:"user_id,omitempty"`
	CreatedAt    *time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy    *string         `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt    *time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy    *string         `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty,omitempty"`
	DeletedBy    *string         `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
