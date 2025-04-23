package assets

type AssetMaintenanceTypeRequest struct {
	MaintenanceTypeName string `json:"maintenance_type_name" gorm:"unique"`
	Description         string `json:"description"`
}
