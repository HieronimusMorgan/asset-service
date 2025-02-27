package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetImage struct {
	ImageID      uint           `gorm:"primaryKey" json:"image_id"`
	UserClientID string         `gorm:"type:varchar(255)" json:"user_client_id"`
	AssetID      uint           `gorm:"not null" json:"asset_id"`
	ImageURL     string         `gorm:"not null" json:"image_url"`
	FileType     string         `gorm:"not null" json:"file_type"`
	FileSize     int64          `gorm:"type:int" json:"file_size"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy    string         `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy    string         `gorm:"type:varchar(255)" json:"updated_by"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy    *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}

type ImageDeleteRequest struct {
	ClientID string   `json:"client_id"`
	Images   []string `json:"images"`
}
