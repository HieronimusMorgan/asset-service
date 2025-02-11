package assets

import (
	"asset-service/internal/models/assets"
	"gorm.io/gorm"
)

type AssetMaintenanceTypeRepository interface {
	GetAssetMaintenanceTypeByName(name string) error
	AddAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType) error
	GetAssetMaintenanceType() ([]assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByID(assetMaintenanceTypeID uint) (*assets.AssetMaintenanceType, error)
	UpdateAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType) error
	DeleteAssetMaintenanceTypeByID(assetMaintenanceTypeID uint) error
	DeleteAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType) error
	GetAssetMaintenanceTypeByType(maintenanceType string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndID(maintenanceType string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeNotExist(maintenanceType string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndIDNotExist(maintenanceType string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndNameNotExist(maintenanceType string, maintenanceTypeName string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndName(maintenanceType string, maintenanceTypeName string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndNameAndID(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndNameAndIDNotExist(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error)
}

type assetMaintenanceTypeRepository struct {
	db gorm.DB
}

const tableAssetMaintenanceTypeName = "my-home.asset_maintenance_type"

func NewAssetMaintenanceTypeRepository(db gorm.DB) AssetMaintenanceTypeRepository {
	return assetMaintenanceTypeRepository{db: db}
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByName(name string) error {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name LIKE ?", name).First(&assetMaintenanceType).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetMaintenanceTypeRepository) AddAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType) error {
	err := r.db.Table(tableAssetMaintenanceTypeName).Create(assetMaintenanceType).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceType() ([]assets.AssetMaintenanceType, error) {
	var assetMaintenanceType []assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Find(&assetMaintenanceType).Where("deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}
	return assetMaintenanceType, nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByID(assetMaintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_id = ?", assetMaintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r assetMaintenanceTypeRepository) UpdateAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType) error {
	err := r.db.Table(tableAssetMaintenanceTypeName).Save(assetMaintenanceType).Error
	if err != nil {
		return err
	}
	return nil
}

// delete by id
func (r assetMaintenanceTypeRepository) DeleteAssetMaintenanceTypeByID(assetMaintenanceTypeID uint) error {
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_id = ?", assetMaintenanceTypeID).Delete(&assets.AssetMaintenanceType{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetMaintenanceTypeRepository) DeleteAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType) error {
	err := r.db.Table(tableAssetMaintenanceTypeName).Delete(assetMaintenanceType).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByType(maintenanceType string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ?", maintenanceType).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndID(maintenanceType string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_id != ?", maintenanceType, maintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndIDNotExist(maintenanceType string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_id != ?", maintenanceType, maintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeNotExist(maintenanceType string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ?", maintenanceType).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndNameNotExist(maintenanceType string, maintenanceTypeName string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_name NOT LIKE ?", maintenanceType, maintenanceTypeName).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndName(maintenanceType string, maintenanceTypeName string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_name LIKE ?", maintenanceType, maintenanceTypeName).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndNameAndID(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_name LIKE ? AND maintenance_type_id != ?", maintenanceType, maintenanceTypeName, maintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndNameAndIDNotExist(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(tableAssetMaintenanceTypeName).Where("maintenance_type_name = ? AND maintenance_type_name LIKE ? AND maintenance_type_id != ?", maintenanceType, maintenanceTypeName, maintenanceTypeID).First(&assetMaintenanceType).Error
	if err != nil {
		return nil, err
	}
	return &assetMaintenanceType, nil
}
