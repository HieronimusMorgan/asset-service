package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetGroupPermission struct {
	PermissionID   uint           `gorm:"primaryKey;column:permission_id" json:"permission_id,omitempty"`
	PermissionName string         `gorm:"size:100;not null;unique" json:"permission_name,omitempty"`
	Description    string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy      string         `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy      string         `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty,omitempty"`
	DeletedBy      *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
