package assets

import (
	request "asset-service/internal/dto/in/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
)

type AssetGroupAssetService interface {
	AddAssetGroupAsset(assetRequest *request.AssetGroupAssetRequest, clientID string) (interface{}, error)
	UpdateAssetGroupAsset(permissionID uint, assetCategoryRequest *request.AssetGroupAssetRequest, clientID string) (interface{}, error)
	GetListAssetGroupAsset(clientID string) (interface{}, error)
	GetAssetGroupAssetById(permissionID uint, clientID string) (interface{}, error)
	DeleteAssetGroupAsset(permissionID uint, clientID string) error
}

type assetGroupAssetService struct {
	AssetGroupAssetRepository repository.AssetGroupAssetRepository
	AssetRepository           repository.AssetRepository
	AssetAuditLogRepository   repository.AssetAuditLogRepository
	Redis                     utils.RedisService
}

func NewAssetGroupAssetService(AssetGroupAssetRepository repository.AssetGroupAssetRepository, AssetRepository repository.AssetRepository, AssetAuditLogRepository repository.AssetAuditLogRepository, redis utils.RedisService) AssetGroupAssetService {
	return &assetGroupAssetService{
		AssetGroupAssetRepository: AssetGroupAssetRepository,
		AssetRepository:           AssetRepository,
		AssetAuditLogRepository:   AssetAuditLogRepository,
		Redis:                     redis,
	}
}

func (s *assetGroupAssetService) AddAssetGroupAsset(assetRequest *request.AssetGroupAssetRequest, clientID string) (interface{}, error) {
	return nil, nil
}

func (s *assetGroupAssetService) UpdateAssetGroupAsset(permissionID uint, assetCategoryRequest *request.AssetGroupAssetRequest, clientID string) (interface{}, error) {
	return nil, nil
}

func (s *assetGroupAssetService) GetListAssetGroupAsset(clientID string) (interface{}, error) {
	return nil, nil
}

func (s *assetGroupAssetService) GetAssetGroupAssetById(permissionID uint, clientID string) (interface{}, error) {
	return nil, nil
}

func (s *assetGroupAssetService) DeleteAssetGroupAsset(permissionID uint, clientID string) error {
	return nil
}
