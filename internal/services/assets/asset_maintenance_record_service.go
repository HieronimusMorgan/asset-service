package assets

import (
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/redis"
	"fmt"
)

type AssetMaintenanceRecordService interface {
	GetMaintenanceRecordByID(maintenanceRecordID uint, clientID string) (interface{}, error)
	GetListMaintenancesRecordByAssetID(assetID uint, clientID string) (interface{}, error)
	GetMaintenancesRecordByAssetIDAndMaintenanceID(assetID, maintenanceID uint, clientID string) (interface{}, error)
	GetMaintenanceRecordByRecordIDAndAssetIDAndMaintenanceID(maintenanceRecordID, assetID, maintenanceID uint, clientID string) (interface{}, error)
}

type assetMaintenanceRecordService struct {
	AssetMaintenanceRepository       repository.AssetMaintenanceRepository
	AssetRepository                  repository.AssetRepository
	AssetMaintenanceRecordRepository repository.AssetMaintenanceRecordRepository
	AssetAuditLogRepository          repository.AssetAuditLogRepository
	Redis                            redis.RedisService
}

func NewAssetMaintenanceRecordService(
	AssetMaintenance repository.AssetMaintenanceRepository,
	assetRepository repository.AssetRepository,
	AssetMaintenanceRecordRecord repository.AssetMaintenanceRecordRepository,
	AssetAuditLogRepository repository.AssetAuditLogRepository,
	RedisService redis.RedisService) AssetMaintenanceRecordService {
	return assetMaintenanceRecordService{
		AssetMaintenanceRepository:       AssetMaintenance,
		AssetRepository:                  assetRepository,
		AssetMaintenanceRecordRepository: AssetMaintenanceRecordRecord,
		AssetAuditLogRepository:          AssetAuditLogRepository,
		Redis:                            RedisService}
}

func (s assetMaintenanceRecordService) GetMaintenanceRecordByID(maintenanceRecordID uint, clientID string) (interface{}, error) {
	maintenance, err := s.AssetMaintenanceRecordRepository.GetMaintenanceRecordByMaintenanceID(maintenanceRecordID, clientID)
	if err != nil {
		return nil, err
	}

	return maintenance, nil
}

func (s assetMaintenanceRecordService) GetListMaintenancesRecordByAssetID(assetID uint, clientID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user from Redis")
	}

	asset, _ := s.AssetRepository.GetAssetResponseByID(data.ClientID, assetID)
	if asset == nil {
		return nil, fmt.Errorf("asset not found")
	}

	assetMaintenanceRecord, _ := s.AssetMaintenanceRecordRepository.GetListMaintenanceRecordByAssetID(asset.AssetID, data.ClientID)
	if assetMaintenanceRecord == nil {
		return nil, fmt.Errorf("asset maintenance not found")
	}

	return struct {
		Asset                  interface{} `json:"asset"`
		AssetMaintenanceRecord interface{} `json:"asset_maintenance_record"`
	}{
		Asset:                  asset,
		AssetMaintenanceRecord: assetMaintenanceRecord,
	}, nil
}

func (s assetMaintenanceRecordService) GetMaintenancesRecordByAssetIDAndMaintenanceID(assetID, maintenanceID uint, clientID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user from Redis")
	}

	asset, _ := s.AssetRepository.GetAssetResponseByID(data.ClientID, assetID)
	if asset == nil {
		return nil, fmt.Errorf("asset not found")
	}

	assetMaintenance, _ := s.AssetMaintenanceRepository.GetMaintenanceByMaintenanceIDAndAssetID(maintenanceID, asset.AssetID, data.ClientID)
	if assetMaintenance == nil {
		return nil, fmt.Errorf("asset maintenance not found")
	}

	assetMaintenanceRecord, _ := s.AssetMaintenanceRecordRepository.GetListMaintenanceRecordByAssetIDAndMaintenanceID(asset.AssetID, assetMaintenance.ID, data.ClientID)
	if assetMaintenanceRecord == nil {
		return nil, fmt.Errorf("asset maintenance record not found")
	}

	return struct {
		Asset                  interface{} `json:"asset"`
		AssetMaintenanceRecord interface{} `json:"asset_maintenance_record"`
	}{
		Asset:                  asset,
		AssetMaintenanceRecord: assetMaintenanceRecord,
	}, nil
}

func (s assetMaintenanceRecordService) GetMaintenanceRecordByRecordIDAndAssetIDAndMaintenanceID(maintenanceRecordID, assetID, maintenanceID uint, clientID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user from Redis")
	}

	asset, _ := s.AssetRepository.GetAssetResponseByID(data.ClientID, assetID)
	if asset == nil {
		return nil, fmt.Errorf("asset not found")
	}

	assetMaintenance, _ := s.AssetMaintenanceRepository.GetMaintenanceByMaintenanceIDAndAssetID(maintenanceID, asset.AssetID, data.ClientID)
	if assetMaintenance == nil {
		return nil, fmt.Errorf("asset maintenance not found")
	}

	maintenanceRecord, _ := s.AssetMaintenanceRecordRepository.GetMaintenanceRecordByRecordIDAndAssetIDAndMaintenanceID(maintenanceRecordID, asset.AssetID, assetMaintenance.ID, data.ClientID)
	if maintenanceRecord == nil {
		return nil, fmt.Errorf("asset maintenance record not found")
	}

	return struct {
		Asset                  interface{} `json:"asset"`
		AssetMaintenanceRecord interface{} `json:"asset_maintenance_record"`
	}{
		Asset:                  asset,
		AssetMaintenanceRecord: maintenanceRecord,
	}, nil
}
