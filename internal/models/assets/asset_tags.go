package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetTag struct {
	TagID       uint           `gorm:"primaryKey" json:"tag_id"`
	TagName     string         `gorm:"type:varchar(255);not null;unique" json:"tag_name"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   string         `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   string         `gorm:"type:varchar(255)" json:"updated_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy   *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
