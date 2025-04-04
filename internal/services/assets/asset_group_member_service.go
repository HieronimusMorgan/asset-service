package assets

import (
	request "asset-service/internal/dto/in/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
)

type AssetGroupMemberService interface {
	AddAssetGroupMember(assetRequest *request.AssetGroupMemberRequest, clientID string) (interface{}, error)
	UpdateAssetGroupMember(permissionID uint, assetCategoryRequest *request.AssetGroupMemberRequest, clientID string) (interface{}, error)
	GetListAssetGroupMember(clientID string) (interface{}, error)
	GetAssetGroupMemberById(permissionID uint, clientID string) (interface{}, error)
	DeleteAssetGroupMember(permissionID uint, clientID string) error
}

type assetGroupMemberService struct {
	AssetGroupMemberRepository repository.AssetGroupMemberRepository
	AssetRepository            repository.AssetRepository
	AssetAuditLogRepository    repository.AssetAuditLogRepository
	Redis                      utils.RedisService
}

func NewAssetGroupMemberService(
	assetCategoryRepository repository.AssetGroupMemberRepository,
	AssetRepository repository.AssetRepository,
	AssetAuditLogRepository repository.AssetAuditLogRepository,
	redis utils.RedisService) AssetGroupMemberService {
	return &assetGroupMemberService{
		AssetGroupMemberRepository: assetCategoryRepository,
		AssetRepository:            AssetRepository,
		AssetAuditLogRepository:    AssetAuditLogRepository,
		Redis:                      redis,
	}
}

func (s *assetGroupMemberService) AddAssetGroupMember(assetRequest *request.AssetGroupMemberRequest, clientID string) (interface{}, error) {
	return nil, nil
}

func (s *assetGroupMemberService) UpdateAssetGroupMember(permissionID uint, assetCategoryRequest *request.AssetGroupMemberRequest, clientID string) (interface{}, error) {
	return nil, nil
}

func (s *assetGroupMemberService) GetListAssetGroupMember(clientID string) (interface{}, error) {
	return nil, nil
}

func (s *assetGroupMemberService) GetAssetGroupMemberById(permissionID uint, clientID string) (interface{}, error) {
	return nil, nil
}

func (s *assetGroupMemberService) DeleteAssetGroupMember(permissionID uint, clientID string) error {
	return nil
}
