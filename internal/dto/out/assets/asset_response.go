package assets

import (
	"time"
)

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if time.Time(d).IsZero() {
		return []byte(`null`), nil
	}
	formatted := time.Time(d).Format(`"2006-01-02"`)
	return []byte(formatted), nil
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = DateOnly(time.Time{})
		return nil
	}
	parsed, err := time.Parse(`"2006-01-02"`, string(data))
	if err != nil {
		return err
	}
	*d = DateOnly(parsed)
	return nil
}

type AssetResponse struct {
	AssetID            uint                  `json:"asset_id,omitempty"`
	UserClientID       string                `json:"user_client_id,omitempty"`
	SerialNumber       *string               `json:"serial_number,omitempty"`
	Name               string                `json:"name,omitempty"`
	Description        string                `json:"description,omitempty"`
	Barcode            *string               `json:"barcode,omitempty"`
	Status             AssetStatusResponse   `json:"status,omitempty"`
	Category           AssetCategoryResponse `json:"category,omitempty"`
	Images             []AssetImageResponse  `json:"images,omitempty"`
	PurchaseDate       *DateOnly             `json:"purchase_date,omitempty"`
	ExpiryDate         *DateOnly             `json:"expiry_date,omitempty"`
	WarrantyExpiryDate *DateOnly             `json:"warranty_expiry_date,omitempty"`
	Price              float64               `json:"price,omitempty"`
	Stock              AssetStockResponse    `json:"stock,omitempty"`
	Notes              *string               `json:"notes,omitempty"`
}

type AssetResponseAll struct {
	ID             uint                  `json:"asset_id,omitempty"`
	ClientID       string                `json:"user_client_id,omitempty"`
	SerialNumber   string                `json:"serial_number,omitempty"`
	Name           string                `json:"name,omitempty"`
	Description    string                `json:"description,omitempty"`
	Barcode        string                `json:"barcode,omitempty"`
	Status         AssetStatusResponse   `json:"status,omitempty"`
	Category       AssetCategoryResponse `json:"category,omitempty"`
	Images         []AssetImageResponse  `json:"images,omitempty"`
	PurchaseDate   string                `json:"purchase_date,omitempty"`
	ExpiryDate     string                `json:"expiry_date,omitempty"`
	WarrantyExpiry string                `json:"warranty_expiry_date,omitempty"`
	Price          float64               `json:"price,omitempty"`
	Stock          int                   `json:"stock,omitempty"`
}

type AssetImageResponse struct {
	ImageURL   string    `json:"image_url,omitempty"`
	FileType   string    `json:"file_type,omitempty"`
	FileSize   int64     `json:"file_size,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty"`
}

type AssetStockResponse struct {
	StockID         uint    `json:"stock_id,omitempty"`
	AssetID         uint    `json:"asset_id,omitempty"`
	InitialQuantity int     `json:"initial_quantity,omitempty"`
	LatestQuantity  int     `json:"latest_quantity,omitempty"`
	ChangeType      string  `json:"change_type,omitempty"`
	Quantity        int     `json:"quantity,omitempty"`
	Reason          *string `json:"reason,omitempty,omitempty"`
}
