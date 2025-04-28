package assets

import (
	"time"
)

// AssetStockHistory represents the history of stock changes for an asset
type AssetStockHistory struct {
	StockHistoryID   uint      `gorm:"primaryKey;autoIncrement" json:"stock_history_id"`
	AssetID          uint      `gorm:"not null;index" json:"asset_id"`
	UserClientID     string    `gorm:"type:varchar(50);not null;index" json:"user_client_id"`
	StockID          uint      `gorm:"not null;index" json:"stock_id"`
	ChangeType       string    `gorm:"type:varchar(50);not null;check:change_type IN ('INCREASE', 'DECREASE', 'ADJUSTMENT')" json:"change_type"`
	PreviousQuantity int       `gorm:"not null;check:previous_quantity >= 0" json:"previous_quantity"`
	NewQuantity      int       `gorm:"not null;check:new_quantity >= 0" json:"new_quantity"`
	QuantityChanged  int       `gorm:"not null;check:quantity_changed > 0" json:"quantity_changed"`
	Reason           *string   `gorm:"type:text" json:"reason,omitempty"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy        *string   `gorm:"type:varchar(255)" json:"created_by"`
	// Relationships
	Asset Asset      `gorm:"foreignKey:AssetID;constraint:OnDelete:CASCADE"`
	Stock AssetStock `gorm:"foreignKey:StockID;constraint:OnDelete:CASCADE"`
}
