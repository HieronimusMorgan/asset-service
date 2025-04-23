package assets

import "time"

type AssetWishlist struct {
	WishlistID    uint       `gorm:"primaryKey;column:wishlist_id" json:"wishlist_id"`
	UserClientID  string     `gorm:"size:50;not null;column:user_client_id" json:"user_client_id"`
	AssetName     string     `gorm:"size:100;not null;column:asset_name" json:"asset_name"`
	SerialNumber  *string    `gorm:"size:100;column:serial_number" json:"serial_number,omitempty"`
	Barcode       *string    `gorm:"size:100;column:barcode" json:"barcode,omitempty"`
	CategoryID    int        `gorm:"column:category_id" json:"category_id,omitempty"`
	StatusID      int        `gorm:"column:status_id" json:"status_id,omitempty"`
	PriorityLevel string     `gorm:"size:20;default:medium;check:priority_level IN ('low','medium','high');column:priority_level" json:"priority_level"`
	PriceEstimate float64    `gorm:"type:decimal(40,2);default:0;column:price_estimate" json:"price_estimate"`
	Notes         *string    `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	CreatedBy     *string    `gorm:"size:255;column:created_by" json:"created_by,omitempty"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	UpdatedBy     *string    `gorm:"size:255;column:updated_by" json:"updated_by,omitempty"`
	DeletedAt     *time.Time `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`
	DeletedBy     *string    `gorm:"size:255;column:deleted_by" json:"deleted_by,omitempty"`
}
