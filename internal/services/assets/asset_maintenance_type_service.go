package assets

import (
	model "asset-service/internal/models/assets"
	repo "asset-service/internal/repository/assets"
	"gorm.io/gorm"
)

type AssetMaintenanceTypeService struct {
	Repo *repo.AssetMaintenanceTypeRepository
}

func NewAssetMaintenanceTypeService(db *gorm.DB) *AssetMaintenanceTypeService {
	return &AssetMaintenanceTypeService{Repo: repo.NewAssetMaintenanceTypeRepository(db)}
}

func (s *AssetMaintenanceTypeService) GetMaintenanceTypeByID(maintenanceTypeID uint) (interface{}, error) {
	maintenanceType, err := s.Repo.GetAssetMaintenanceTypeByID(maintenanceTypeID)
	if err != nil {
		return nil, err
	}

	if maintenanceType.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return maintenanceType, nil
}

func (s *AssetMaintenanceTypeService) GetMaintenanceType() ([]model.AssetMaintenanceType, error) {
	maintenanceTypes, err := s.Repo.GetAssetMaintenanceType()
	if err != nil {
		return nil, err
	}

	return maintenanceTypes, nil
}

func (s *AssetMaintenanceTypeService) AddMaintenanceType(maintenanceType *model.AssetMaintenanceType) (interface{}, error) {
	err := s.Repo.AddAssetMaintenanceType(maintenanceType)
	if err != nil {
		return nil, err
	}

	return maintenanceType, nil
}

func (s *AssetMaintenanceTypeService) UpdateMaintenanceType(maintenanceType *model.AssetMaintenanceType) (interface{}, error) {
	err := s.Repo.UpdateAssetMaintenanceType(maintenanceType)
	if err != nil {
		return nil, err
	}

	return maintenanceType, nil
}

func (s *AssetMaintenanceTypeService) DeleteMaintenanceType(maintenanceTypeID uint) error {
	err := s.Repo.DeleteAssetMaintenanceTypeByID(maintenanceTypeID)
	if err != nil {
		return err
	}

	return nil
}
