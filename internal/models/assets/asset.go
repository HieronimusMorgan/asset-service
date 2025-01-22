package assets

import (
	"gorm.io/gorm"
	"time"
)

type Asset struct {
	AssetID      uint           `gorm:"primaryKey" json:"asset_id"`
	UserClientID string         `gorm:"type:varchar(50);not null" json:"user_client_id"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	CategoryID   int            `gorm:"not null" json:"category_id"`
	StatusID     int            `gorm:"not null" json:"status_id"`
	PurchaseDate *time.Time     `gorm:"type:date" json:"purchase_date"`
	ExpiryDate   *time.Time     `gorm:"type:date" json:"expiry_date"`
	Value        float64        `gorm:"type:decimal(15,2)" json:"value"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy    string         `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy    string         `gorm:"type:varchar(50)" json:"updated_by"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy    *string        `gorm:"type:varchar(50)" json:"deleted_by,omitempty"`
}
