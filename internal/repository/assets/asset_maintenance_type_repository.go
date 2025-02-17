package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"gorm.io/gorm"
)

// AssetMaintenanceTypeRepository defines the interface for managing asset maintenance types
type AssetMaintenanceTypeRepository interface {
	GetAssetMaintenanceTypeByName(name string, clientID string) (*assets.AssetMaintenanceType, error)
	AddAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType, clientID string) error
	GetAssetMaintenanceType(clientID string) ([]assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByID(assetMaintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error)
	UpdateAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType, clientID string) error
	DeleteAssetMaintenanceTypeByID(assetMaintenanceTypeID uint, clientID string) error
	DeleteAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType, clientID string) error
	GetAssetMaintenanceTypeByType(maintenanceType string, clientID string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndID(maintenanceType string, maintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeNotExist(maintenanceType string, clientID string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndIDNotExist(maintenanceType string, maintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndNameNotExist(maintenanceType string, maintenanceTypeName string, clientID string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndName(maintenanceType string, maintenanceTypeName string, clientID string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndNameAndID(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error)
	GetAssetMaintenanceTypeByTypeAndNameAndIDNotExist(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error)
}

type assetMaintenanceTypeRepository struct {
	db gorm.DB
}

func NewAssetMaintenanceTypeRepository(db gorm.DB) AssetMaintenanceTypeRepository {
	return assetMaintenanceTypeRepository{db: db}
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByName(name string, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_name LIKE ? AND user_client_id = ?", name, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}

func (r assetMaintenanceTypeRepository) AddAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType, clientID string) error {
	assetMaintenanceType.UserClientID = clientID
	return r.db.Table(utils.TableAssetMaintenanceTypeName).Create(assetMaintenanceType).Error
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceType(clientID string) ([]assets.AssetMaintenanceType, error) {
	var assetMaintenanceTypes []assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("user_client_id = ? AND deleted_at IS NULL", clientID).
		Find(&assetMaintenanceTypes).Error
	return assetMaintenanceTypes, err
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByID(assetMaintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_id = ? AND user_client_id = ?", assetMaintenanceTypeID, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}

func (r assetMaintenanceTypeRepository) UpdateAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType, clientID string) error {
	assetMaintenanceType.UserClientID = clientID
	return r.db.Table(utils.TableAssetMaintenanceTypeName).Save(assetMaintenanceType).Error
}

func (r assetMaintenanceTypeRepository) DeleteAssetMaintenanceTypeByID(assetMaintenanceTypeID uint, clientID string) error {
	return r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_id = ? AND user_client_id = ?", assetMaintenanceTypeID, clientID).
		Delete(&assets.AssetMaintenanceType{}).Error
}

func (r assetMaintenanceTypeRepository) DeleteAssetMaintenanceType(assetMaintenanceType *assets.AssetMaintenanceType, clientID string) error {
	return r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("user_client_id = ?", clientID).
		Delete(assetMaintenanceType).Error
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByType(maintenanceType string, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_name = ? AND user_client_id = ?", maintenanceType, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndID(maintenanceType string, maintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_name = ? AND maintenance_type_id != ? AND user_client_id = ?", maintenanceType, maintenanceTypeID, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}

// Additional filters with clientID
func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeNotExist(maintenanceType string, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_name = ? AND user_client_id = ?", maintenanceType, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndIDNotExist(maintenanceType string, maintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_name = ? AND maintenance_type_id != ? AND user_client_id = ?", maintenanceType, maintenanceTypeID, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndNameNotExist(maintenanceType string, maintenanceTypeName string, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_name = ? AND maintenance_type_name NOT LIKE ? AND user_client_id = ?", maintenanceType, maintenanceTypeName, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndName(maintenanceType string, maintenanceTypeName string, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_name = ? AND maintenance_type_name LIKE ? AND user_client_id = ?", maintenanceType, maintenanceTypeName, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndNameAndID(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_name = ? AND maintenance_type_name LIKE ? AND maintenance_type_id = ? AND user_client_id = ?", maintenanceType, maintenanceTypeName, maintenanceTypeID, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}

func (r assetMaintenanceTypeRepository) GetAssetMaintenanceTypeByTypeAndNameAndIDNotExist(maintenanceType string, maintenanceTypeName string, maintenanceTypeID uint, clientID string) (*assets.AssetMaintenanceType, error) {
	var assetMaintenanceType assets.AssetMaintenanceType
	err := r.db.Table(utils.TableAssetMaintenanceTypeName).
		Where("maintenance_type_name = ? AND maintenance_type_name NOT LIKE ? AND maintenance_type_id != ? AND user_client_id = ?", maintenanceType, maintenanceTypeName, maintenanceTypeID, clientID).
		First(&assetMaintenanceType).Error
	return &assetMaintenanceType, err
}
