package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetGroupInvitation struct {
	InvitationID     uint           `gorm:"primaryKey" json:"invitation_id,omitempty"`
	AssetGroupID     uint           `json:"asset_group_id,omitempty"`
	InvitedUserID    uint           `json:"invited_user_id,omitempty"`
	InvitedUserToken string         `gorm:"unique" json:"invited_user_token,omitempty"`
	InvitedByUserID  uint           `json:"invited_by_user_id,omitempty"`
	Status           string         `json:"status,omitempty"`
	Message          string         `json:"message,omitempty"`
	InvitedAt        string         `json:"invited_at,omitempty"`
	RespondedAt      string         `json:"responded_at,omitempty"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy        string         `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy        string         `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty,omitempty"`
	DeletedBy        *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
