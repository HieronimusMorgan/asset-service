package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
)

type AssetGroupPermissionService interface {
	AddAssetGroupPermission(assetGroupPermissionRequest *request.AssetGroupPermissionRequest, clientID string) error
	UpdateAssetGroupPermission(permissionID uint, assetCategoryRequest *request.AssetGroupPermissionRequest, clientID string) error
	GetListAssetGroupPermission(clientID string) (interface{}, error)
	GetAssetGroupPermissionById(permissionID uint, clientID string) (interface{}, error)
	DeleteAssetGroupPermission(permissionID uint, clientID string) error
}

type assetGroupPermissionService struct {
	AssetGroupPermissionRepository repository.AssetGroupPermissionRepository
	AssetRepository                repository.AssetRepository
	AssetAuditLogRepository        repository.AssetAuditLogRepository
	Redis                          utils.RedisService
}

func NewAssetGroupPermissionService(
	assetCategoryRepository repository.AssetGroupPermissionRepository,
	AssetRepository repository.AssetRepository,
	AssetAuditLogRepository repository.AssetAuditLogRepository,
	redis utils.RedisService) AssetGroupPermissionService {
	return &assetGroupPermissionService{
		AssetGroupPermissionRepository: assetCategoryRepository,
		AssetRepository:                AssetRepository,
		AssetAuditLogRepository:        AssetAuditLogRepository,
		Redis:                          redis,
	}
}

func (s *assetGroupPermissionService) AddAssetGroupPermission(req *request.AssetGroupPermissionRequest, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	assetPermission := &assets.AssetGroupPermission{
		PermissionName: req.PermissionName,
		Description:    req.Description,
		CreatedBy:      data.ClientID,
	}

	if err := s.AssetGroupPermissionRepository.AddAssetGroupPermission(assetPermission); err != nil {
		return logErrorWithNoReturn("AddAssetGroupPermission", clientID, err, "Failed to add asset group permission")
	}

	return nil
}

func (s *assetGroupPermissionService) UpdateAssetGroupPermission(permissionID uint, assetCategoryRequest *request.AssetGroupPermissionRequest, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	assetPermission, err := s.AssetGroupPermissionRepository.GetAssetGroupPermissionByID(permissionID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupPermissionById", clientID, err, "Failed to get asset group permission")
	}

	assetPermission.PermissionName = assetCategoryRequest.PermissionName
	assetPermission.Description = assetCategoryRequest.Description
	assetPermission.UpdatedBy = data.ClientID

	if err := s.AssetGroupPermissionRepository.UpdateAssetGroupPermission(assetPermission); err != nil {
		return logErrorWithNoReturn("UpdateAssetGroupPermission", clientID, err, "Failed to update asset group permission")
	}

	return nil
}

func (s *assetGroupPermissionService) GetListAssetGroupPermission(clientID string) (interface{}, error) {
	_, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	permissions, err := s.AssetGroupPermissionRepository.GetListAssetGroupPermission()
	if err != nil {
		return logError("GetListAssetGroupPermission", clientID, err, "Failed to get list of asset group permissions")
	}

	return permissions, err
}

func (s *assetGroupPermissionService) GetAssetGroupPermissionById(permissionID uint, clientID string) (interface{}, error) {
	_, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	permission, err := s.AssetGroupPermissionRepository.GetAssetGroupPermissionByID(permissionID)
	if err != nil {
		return logError("GetAssetGroupPermissionByID", clientID, err, "Failed to get asset group permission by MaintenanceTypeID")
	}

	return permission, nil
}

func (s *assetGroupPermissionService) DeleteAssetGroupPermission(permissionID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	assetPermission, err := s.AssetGroupPermissionRepository.GetAssetGroupPermissionByID(permissionID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupPermissionById", clientID, err, "Failed to get asset group permission")
	}

	assetPermission.DeletedBy = &data.ClientID

	if err := s.AssetGroupPermissionRepository.DeleteAssetGroupPermission(assetPermission); err != nil {
		return logErrorWithNoReturn("DeleteAssetGroupPermission", clientID, err, "Failed to delete asset group permission")
	}

	return nil
}
