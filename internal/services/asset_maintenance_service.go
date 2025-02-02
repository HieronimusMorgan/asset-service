package services

import (
	"asset-service/internal/dto/in"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	"asset-service/internal/repository"
	"asset-service/internal/utils"
	"gorm.io/gorm"
)

type AssetMaintenanceService struct {
	Repo *repository.AssetMaintenanceRepository
}

func NewAssetMaintenanceService(db *gorm.DB) *AssetMaintenanceService {
	return &AssetMaintenanceService{Repo: repository.NewAssetMaintenanceRepository(db)}
}

func (s *AssetMaintenanceService) CreateMaintenance(maintenance in.AssetMaintenanceRequest, clientID string) (assets.AssetMaintenance, error) {
	user := &user.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return assets.AssetMaintenance{}, err
	}
	maintenanceRecord := assets.AssetMaintenance{
		AssetID:            maintenance.AssetID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    maintenance.MaintenanceDate,
		CreatedBy:          user.FullName,
		UpdatedBy:          user.FullName,
	}

	if maintenance.MaintenanceDetails == nil {
		maintenanceRecord.MaintenanceDetails = nil
	}

	err = s.Repo.Create(clientID, &maintenanceRecord)
	if err != nil {
		return assets.AssetMaintenance{}, err
	}

	return maintenanceRecord, nil
}

func (s *AssetMaintenanceService) GetMaintenanceByID(maintenanceID uint, clientID string) (interface{}, error) {
	maintenance, err := s.Repo.GetByID(maintenanceID, clientID)
	if err != nil {
		return nil, err
	}

	if maintenance.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return maintenance, nil
}

func (s *AssetMaintenanceService) UpdateMaintenance(maintenance *in.AssetMaintenanceRequest) error {
	return nil
}

func (s *AssetMaintenanceService) DeleteMaintenance(maintenanceID uint) error {
	return s.Repo.Delete(maintenanceID)
}

func (s *AssetMaintenanceService) GetMaintenancesByAssetID(assetID uint, clientID string) (*assets.AssetMaintenance, error) {
	return s.Repo.GetByAssetID(assetID, clientID)
}
