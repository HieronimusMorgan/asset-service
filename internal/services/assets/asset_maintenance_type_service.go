package assets

import (
	model "asset-service/internal/models/assets"
	repo "asset-service/internal/repository/assets"
	"gorm.io/gorm"
)

type AssetMaintenanceTypeService interface {
	GetMaintenanceTypeByID(maintenanceTypeID uint) (interface{}, error)
	GetMaintenanceType() ([]model.AssetMaintenanceType, error)
	AddMaintenanceType(maintenanceType *model.AssetMaintenanceType) (interface{}, error)
	UpdateMaintenanceType(maintenanceType *model.AssetMaintenanceType) (interface{}, error)
	DeleteMaintenanceType(maintenanceTypeID uint) error
}
type assetMaintenanceTypeService struct {
	AssetMaintenanceTypeRepository repo.AssetMaintenanceTypeRepository
}

func NewAssetMaintenanceTypeService(assetMaintenanceTypeRepository repo.AssetMaintenanceTypeRepository) AssetMaintenanceTypeService {
	return assetMaintenanceTypeService{AssetMaintenanceTypeRepository: assetMaintenanceTypeRepository}
}

func (s assetMaintenanceTypeService) GetMaintenanceTypeByID(maintenanceTypeID uint) (interface{}, error) {
	maintenanceType, err := s.AssetMaintenanceTypeRepository.GetAssetMaintenanceTypeByID(maintenanceTypeID)
	if err != nil {
		return nil, err
	}

	if maintenanceType.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return maintenanceType, nil
}

func (s assetMaintenanceTypeService) GetMaintenanceType() ([]model.AssetMaintenanceType, error) {
	maintenanceTypes, err := s.AssetMaintenanceTypeRepository.GetAssetMaintenanceType()
	if err != nil {
		return nil, err
	}

	return maintenanceTypes, nil
}

func (s assetMaintenanceTypeService) AddMaintenanceType(maintenanceType *model.AssetMaintenanceType) (interface{}, error) {
	err := s.AssetMaintenanceTypeRepository.AddAssetMaintenanceType(maintenanceType)
	if err != nil {
		return nil, err
	}

	return maintenanceType, nil
}

func (s assetMaintenanceTypeService) UpdateMaintenanceType(maintenanceType *model.AssetMaintenanceType) (interface{}, error) {
	err := s.AssetMaintenanceTypeRepository.UpdateAssetMaintenanceType(maintenanceType)
	if err != nil {
		return nil, err
	}

	return maintenanceType, nil
}

func (s assetMaintenanceTypeService) DeleteMaintenanceType(maintenanceTypeID uint) error {
	err := s.AssetMaintenanceTypeRepository.DeleteAssetMaintenanceTypeByID(maintenanceTypeID)
	if err != nil {
		return err
	}

	return nil
}
