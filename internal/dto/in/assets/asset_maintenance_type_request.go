package assets

type AssetMaintenanceTypeRequest struct {
	TypeName    string `json:"type_name" gorm:"unique"`
	Description string `json:"description"`
}
