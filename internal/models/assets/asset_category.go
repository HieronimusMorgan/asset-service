package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetCategory struct {
	AssetCategoryID uint           `gorm:"primaryKey" json:"asset_category_id"`
	UserClientID    string         `gorm:"type:varchar(50);not null" json:"user_client_id,omitempty"`
	CategoryName    string         `gorm:"type:varchar(255);not null" json:"category_name"`
	Description     string         `gorm:"type:text" json:"description"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       string         `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       string         `gorm:"type:varchar(255)" json:"updated_by"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy       *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
