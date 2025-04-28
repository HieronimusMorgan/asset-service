package assets

import "mime/multipart"

type AssetRequest struct {
	SerialNumber   *string                 `json:"serial_number"`
	Name           string                  `json:"name"`
	Description    *string                 `json:"description"`
	Barcode        *string                 `json:"barcode"`
	Images         []*multipart.FileHeader `form:"images"` // List of image files
	CategoryID     uint                    `json:"category_id"`
	StatusID       uint                    `json:"status_id"`
	PurchaseDate   *string                 `json:"purchase_date"`
	ExpiryDate     *string                 `json:"expiry_date"`
	WarrantyExpiry *string                 `json:"warranty_expiry_date,omitempty"`
	Price          float64                 `json:"price"`
	Stock          int                     `json:"stock"`
	Notes          *string                 `json:"notes"`
}

func (a *AssetRequest) ConvertAssetRequestEmptyToNil() {
	a.SerialNumber = checkEmptyString(a.SerialNumber)
	a.Description = checkEmptyString(a.Description)
	a.Barcode = checkEmptyString(a.Barcode)
	a.PurchaseDate = checkEmptyString(a.PurchaseDate)
	a.ExpiryDate = checkEmptyString(a.ExpiryDate)
	a.WarrantyExpiry = checkEmptyString(a.WarrantyExpiry)
	a.Notes = checkEmptyString(a.Notes)
}

func checkEmptyString(field *string) *string {
	if field != nil && *field == "" {
		return nil
	}
	return field
}

type UpdateAssetRequest struct {
	SerialNumber       *string `json:"serial_number,omitempty"`
	Description        *string `json:"description,omitempty"`
	Barcode            *string `json:"barcode,omitempty"`
	CategoryID         uint    `json:"category_id,omitempty"`
	StatusID           uint    `json:"status_id,omitempty"`
	PurchaseDate       *string `json:"purchase_date,omitempty"`
	ExpiryDate         *string `json:"expiry_date,omitempty"`
	WarrantyExpiryDate *string `json:"warranty_expiry_date,omitempty"`
	Price              float64 `json:"price,omitempty"`
	Stock              int     `json:"stock,omitempty"`
	Notes              *string `json:"notes"`
}

type AssetWishlistRequest struct {
	AssetName     string  `json:"asset_name"`
	SerialNumber  *string `json:"serial_number,omitempty"`
	Barcode       *string `json:"barcode,omitempty"`
	CategoryID    uint    `json:"category_id,omitempty"`
	StatusID      uint    `json:"status_id,omitempty"`
	PriorityLevel string  `json:"priority_level"`
	PriceEstimate float64 `json:"price_estimate"`
	Notes         *string `json:"notes,omitempty"`
}

type UpdateAssetWishlistRequest struct {
	Description  *string `json:"description"`
	CategoryID   int     `json:"category_id" binding:"required"`
	StatusID     int     `json:"status_id" binding:"required"`
	PurchaseDate *string `json:"purchase_date"`
	Price        float64 `json:"price" binding:"required"`
	Notes        *string `json:"notes"`
}

type UpdateAssetStockRequest struct {
	AssetID  uint    `json:"asset_id"`
	Quantity int     `json:"quantity"`
	Reason   *string `json:"reason,omitempty"`
}
