package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetTagMap struct {
	AssetID   uint           `gorm:"not null" json:"asset_id"`
	TagID     uint           `gorm:"not null" json:"tag_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(255)" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
