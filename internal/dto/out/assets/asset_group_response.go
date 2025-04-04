package assets

type AssetGroupResponse struct {
	AssetGroupID   uint   `gorm:"primaryKey;column:asset_group_id"  json:"asset_group_id,omitempty"`
	AssetGroupName string `gorm:"type:varchar(100);not null"  json:"asset_group_name,omitempty"`
	Description    string `gorm:"type:text" json:"description,omitempty"`
}

type AssetGroupDetailResponse struct {
	AssetGroupID   uint                       `gorm:"primaryKey;column:asset_group_id"  json:"asset_group_id,omitempty"`
	AssetGroupName string                     `gorm:"type:varchar(100);not null"  json:"asset_group_name,omitempty"`
	Description    string                     `gorm:"type:text" json:"description,omitempty"`
	OwnerUserID    uint                       `json:"owner_user_id,omitempty"`
	OwnerName      string                     `json:"owner_name,omitempty"`
	Member         []AssetGroupMemberResponse `json:"member,omitempty"`
}

type AssetGroupMemberResponse struct {
	UserID         uint                                     `json:"user_id"`
	Username       string                                   `json:"username"`
	FullName       string                                   `json:"full_name"`
	ProfilePicture string                                   `json:"profile_picture"`
	Permission     []AssetGroupMemberWithPermissionResponse `json:"permission"`
}

type AssetGroupMemberWithPermissionResponse struct {
	PermissionID   *uint   `json:"permission_id"`
	PermissionName *string `json:"permission_name"`
}

type AssetGroupAssetResponse struct {
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
