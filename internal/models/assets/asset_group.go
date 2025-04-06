package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetGroup struct {
	AssetGroupID    uint           `gorm:"primaryKey;column:asset_group_id"  json:"asset_group_id,omitempty"`
	AssetGroupName  string         `gorm:"type:varchar(100);not null"  json:"asset_group_name,omitempty"`
	Description     string         `gorm:"type:text" json:"description,omitempty"`
	OwnerUserID     uint           `gorm:"column:owner_user_id"  json:"owner_user_id,omitempty"`
	InvitationToken *string        `gorm:"type:varchar(100);unique" json:"invitation_token,omitempty"`
	MaxUses         *int           `gorm:"default:null" json:"max_uses,omitempty"`
	CurrentUses     *int           `gorm:"default:0" json:"current_uses,omitempty"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy       string         `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy       string         `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty,omitempty"`
	DeletedBy       *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
