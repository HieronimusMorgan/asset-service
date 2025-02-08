package assets

import (
	assets3 "asset-service/internal/dto/in/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	assets2 "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"encoding/json"
	"gorm.io/gorm"
	"log"
)

type AssetMaintenanceService struct {
	Repo *assets2.AssetMaintenanceRepository
}

func NewAssetMaintenanceService(db *gorm.DB) *AssetMaintenanceService {
	return &AssetMaintenanceService{Repo: assets2.NewAssetMaintenanceRepository(db)}
}

func (s *AssetMaintenanceService) CreateMaintenance(maintenance assets3.AssetMaintenanceRequest, clientID string) (assets.AssetMaintenance, error) {
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

func (s *AssetMaintenanceService) UpdateMaintenance(maintenance *assets3.AssetMaintenanceRequest) error {
	return nil
}

func (s *AssetMaintenanceService) DeleteMaintenance(maintenanceID uint) error {
	return s.Repo.Delete(maintenanceID)
}

func (s *AssetMaintenanceService) GetMaintenancesByAssetID(assetID uint, clientID string) (*assets.AssetMaintenance, error) {
	return s.Repo.GetByAssetID(assetID, clientID)
}

func (s *AssetMaintenanceService) PerformMaintenanceCheck() error {
	maintenance, err := s.Repo.GetList()
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
