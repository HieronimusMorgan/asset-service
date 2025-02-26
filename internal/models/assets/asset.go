package assets

import (
	"gorm.io/gorm"
	"time"
)

type Asset struct {
	AssetID            uint           `gorm:"primaryKey" json:"asset_id,omitempty"`
	UserClientID       string         `gorm:"type:varchar(50);not null" json:"user_client_id,omitempty"`
	SerialNumber       *string        `gorm:"type:varchar(100)" json:"serial_number,omitempty"`
	Name               string         `gorm:"type:varchar(100);not null" json:"name,omitempty"`
	Description        *string        `gorm:"type:text" json:"description,omitempty"`
	Barcode            *string        `gorm:"type:varchar(100)" json:"barcode,omitempty"`
	CategoryID         int            `gorm:"not null" json:"category_id,omitempty"`
	StatusID           int            `gorm:"not null" json:"status_id,omitempty"`
	PurchaseDate       *time.Time     `gorm:"type:date" json:"purchase_date,omitempty"`
	ExpiryDate         *time.Time     `gorm:"type:date" json:"expiry_date,omitempty"`
	WarrantyExpiryDate *time.Time     `gorm:"type:date" json:"warranty_expiry_date,omitempty"`
	Price              float64        `gorm:"type:decimal(15,2)" json:"price,omitempty"`
	Stock              int            `gorm:"not null" json:"stock,omitempty"`
	Notes              *string        `gorm:"type:text" json:"notes,omitempty"`
	IsWishlist         bool           `gorm:"type:boolean" json:"is_wishlist,omitempty"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy          string         `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy          string         `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty,omitempty"`
	DeletedBy          *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}
