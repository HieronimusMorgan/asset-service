package assets

type AssetRequest struct {
	AssetCode      string  `json:"asset_code"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Barcode        string  `json:"barcode"`
	CategoryID     int     `json:"category_id"`
	StatusID       int     `json:"status_id"`
	PurchaseDate   string  `json:"purchase_date"`
	ExpiryDate     string  `json:"expiry_date"`
	WarrantyExpiry string  `json:"warranty_expiry_date,omitempty"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
}

type AssetWishlistRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	CategoryID   int     `json:"category_id" binding:"required"`
	StatusID     int     `json:"status_id" binding:"required"`
	PurchaseDate string  `json:"purchase_date"`
	Price        float64 `json:"price" binding:"required"`
	IsWishlist   bool    `json:"is_wishlist" binding:"required"`
}
