package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"encoding/json"
	"gorm.io/gorm"
	"log"
)

type AssetMaintenanceService interface {
	CreateMaintenance(maintenance request.AssetMaintenanceRequest, clientID string) (assets.AssetMaintenance, error)
	GetMaintenanceByID(maintenanceID uint, clientID string) (interface{}, error)
	UpdateMaintenance(maintenance *request.AssetMaintenanceRequest) error
	DeleteMaintenance(maintenanceID uint) error
	GetMaintenancesByAssetID(assetID uint, clientID string) (*assets.AssetMaintenance, error)
	PerformMaintenanceCheck() error
}

type assetMaintenanceService struct {
	AssetMaintenanceRepository repository.AssetMaintenanceRepository
	RedisService               utils.RedisService
}

func NewAssetMaintenanceService(AssetMaintenance repository.AssetMaintenanceRepository, RedisService utils.RedisService) AssetMaintenanceService {
	return assetMaintenanceService{AssetMaintenanceRepository: AssetMaintenance, RedisService: RedisService}
}

func (s assetMaintenanceService) CreateMaintenance(maintenance request.AssetMaintenanceRequest, clientID string) (assets.AssetMaintenance, error) {
	user := &user.User{}
	err := s.RedisService.GetData(utils.User, clientID, user)
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

	err = s.AssetMaintenanceRepository.Create(clientID, &maintenanceRecord)
	if err != nil {
		return assets.AssetMaintenance{}, err
	}

	return maintenanceRecord, nil
}

func (s assetMaintenanceService) GetMaintenanceByID(maintenanceID uint, clientID string) (interface{}, error) {
	maintenance, err := s.AssetMaintenanceRepository.GetByID(maintenanceID, clientID)
	if err != nil {
		return nil, err
	}

	if maintenance.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return maintenance, nil
}

func (s assetMaintenanceService) UpdateMaintenance(maintenance *request.AssetMaintenanceRequest) error {
	return nil
}

func (s assetMaintenanceService) DeleteMaintenance(maintenanceID uint) error {
	return s.AssetMaintenanceRepository.Delete(maintenanceID)
}

func (s assetMaintenanceService) GetMaintenancesByAssetID(assetID uint, clientID string) (*assets.AssetMaintenance, error) {
	return s.AssetMaintenanceRepository.GetByAssetID(assetID, clientID)
}

func (s assetMaintenanceService) PerformMaintenanceCheck() error {
	maintenance, err := s.AssetMaintenanceRepository.GetList()
	if err != nil {
		return err
	}

	for _, m := range maintenance {
		if m.MaintenanceDate == "" {
			continue
		}

		jsonPretty, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return err
		}
		log.Printf("Asset JSON (pretty):\n%s", jsonPretty)
	}
	return nil
}
