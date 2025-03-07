package assets

import (
	"gorm.io/gorm"
	"time"
)

type AssetStock struct {
	StockID         uint           `gorm:"primaryKey" json:"stock_id"`
	AssetID         uint           `gorm:"not null;index" json:"asset_id"`
	UserClientID    string         `gorm:"type:varchar(50);not null" json:"user_client_id,omitempty"`
	InitialQuantity int            `gorm:"not null;check:initial_quantity >= 0" json:"initial_quantity"`
	LatestQuantity  int            `gorm:"not null;check:latest_quantity >= 0" json:"latest_quantity"`
	ChangeType      string         `gorm:"type:varchar(50);not null;check:change_type IN ('INCREASE', 'DECREASE')" json:"change_type"`
	Quantity        int            `gorm:"not null;check:quantity > 0" json:"quantity"`
	Reason          *string        `gorm:"type:text" json:"reason,omitempty"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       string         `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       string         `gorm:"type:varchar(255)" json:"updated_by"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy       *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}

func (AssetStock) TableName() string {
	return "asset_stock"
}
