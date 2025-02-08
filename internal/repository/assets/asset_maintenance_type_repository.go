package assets

import (
	"asset-service/internal/models/assets"
	"gorm.io/gorm"
)

type AssetMaintenanceTypeRepository struct {
	DB *gorm.DB
}

const tableAssetMaintenanceTypeName = "my-home.asset_maintenance_type"

func NewAssetMaintenanceTypeRepository(db *gorm.DB) *AssetMaintenanceTypeRepository {
	return &AssetMaintenanceTypeRepository{DB: db}
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByName(name string) error {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name LIKE ?", name).First(&assetMaintenanceType).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetMaintenanceTypeRepository) AddAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType) error {
	err := r.DB.Table(tableAssetMaintenanceTypeName).Create(assetMaintenanceType).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceType() ([]assets.AssetMaintenanceType, error) {
	var assetMaintenanceType []assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Find(&assetMaintenanceType).Where("deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}
	return assetMaintenanceType, nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByID(assetMaintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_id = ?", assetMaintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r AssetMaintenanceTypeRepository) UpdateAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType) error {
	err := r.DB.Table(tableAssetMaintenanceTypeName).Save(assetMaintenanceType).Error
	if err != nil {
		return err
	}
	return nil
}

// delete by id
func (r AssetMaintenanceTypeRepository) DeleteAssetMaintenanceTypeByID(assetMaintenanceTypeID uint) error {
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_id = ?", assetMaintenanceTypeID).Delete(&assets.AssetMaintenanceType{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetMaintenanceTypeRepository) DeleteAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType) error {
	err := r.DB.Table(tableAssetMaintenanceTypeName).Delete(assetMaintenanceType).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByType(maintenanceType string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ?", maintenanceType).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndID(maintenanceType string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_id != ?", maintenanceType, maintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndIDNotExist(maintenanceType string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_id != ?", maintenanceType, maintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeNotExist(maintenanceType string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ?", maintenanceType).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndNameNotExist(maintenanceType string, maintenanceTypeName string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_name NOT LIKE ?", maintenanceType, maintenanceTypeName).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndName(maintenanceType string, maintenanceTypeName string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_name LIKE ?", maintenanceType, maintenanceTypeName).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndNameAndID(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_name LIKE ? AND maintenance_type_id != ?", maintenanceType, maintenanceTypeName, maintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r AssetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndNameAndIDNotExist(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.DB.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_name LIKE ? AND maintenance_type_id != ?", maintenanceType, maintenanceTypeName, maintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}
