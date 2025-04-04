package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	users "asset-service/internal/repository/users"
	"asset-service/internal/utils"
)

type AssetGroupService interface {
	AddAssetGroup(assetRequest *request.AssetGroupRequest, clientID string) error
	UpdateAssetGroup(assetGroupID uint, assetCategoryRequest *request.AssetGroupRequest, clientID string) (interface{}, error)
	GetListAssetGroup(clientID string) (interface{}, error)
	GetAssetGroupDetailByID(assetGroupID uint, clientID string) (interface{}, error)
	DeleteAssetGroup(assetGroupID uint, clientID string) error
	AddMemberAssetGroup(req *request.AssetGroupMemberRequest, clientID string) error
}

type assetGroupService struct {
	UserRepository             users.UserRepository
	AssetGroupRepository       repository.AssetGroupRepository
	permissionRepository       repository.AssetGroupPermissionRepository
	memberPermissionRepository repository.AssetGroupMemberPermissionRepository
	memberRepository           repository.AssetGroupMemberRepository
	assetGroupAssetRepository  repository.AssetGroupAssetRepository
	AssetRepository            repository.AssetRepository
	AssetAuditLogRepository    repository.AssetAuditLogRepository
	Redis                      utils.RedisService
}

func NewAssetGroupService(UserRepository users.UserRepository, AssetGroupRepository repository.AssetGroupRepository, permissionRepository repository.AssetGroupPermissionRepository, memberPermissionRepository repository.AssetGroupMemberPermissionRepository, memberRepository repository.AssetGroupMemberRepository, assetGroupAssetRepository repository.AssetGroupAssetRepository, AssetRepository repository.AssetRepository, AssetAuditLogRepository repository.AssetAuditLogRepository, redis utils.RedisService) AssetGroupService {
	return &assetGroupService{
		UserRepository:             UserRepository,
		AssetGroupRepository:       AssetGroupRepository,
		permissionRepository:       permissionRepository,
		memberPermissionRepository: memberPermissionRepository,
		memberRepository:           memberRepository,
		assetGroupAssetRepository:  assetGroupAssetRepository,
		AssetRepository:            AssetRepository,
		AssetAuditLogRepository:    AssetAuditLogRepository,
		Redis:                      redis,
	}
}

func (s *assetGroupService) AddAssetGroup(assetRequest *request.AssetGroupRequest, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	// Add asset group
	assetGroup := &assets.AssetGroup{
		GroupName:   assetRequest.GroupName,
		Description: assetRequest.Description,
		OwnerUserID: user.UserID,
		CreatedBy:   user.ClientID,
		UpdatedBy:   user.ClientID,
	}

	err = s.AssetGroupRepository.AddAssetGroup(assetGroup, clientID, user)
	if err != nil {
		return logErrorWithNoReturn("AddAssetGroup", clientID, err, "Failed to add asset group")
	}
	return nil
}

func (s *assetGroupService) UpdateAssetGroup(assetGroupID uint, assetCategoryRequest *request.AssetGroupRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}
	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(assetGroupID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetailByID", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetailByID", clientID, nil, "Asset group not found")
	}

	// Update asset group
	assetGroup.GroupName = assetCategoryRequest.GroupName
	assetGroup.Description = assetCategoryRequest.Description
	assetGroup.UpdatedBy = user.ClientID
	err = s.AssetGroupRepository.UpdateAssetGroup(assetGroup, user.UserID)
	if err != nil {
		return nil, logErrorWithNoReturn("UpdateAssetGroup", clientID, err, "Failed to update asset group")
	}

	return assetGroup, nil
}

func (s *assetGroupService) GetListAssetGroup(clientID string) (interface{}, error) {
	return nil, nil
}

func (s *assetGroupService) GetAssetGroupDetailByID(assetGroupID uint, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupDetailByID(assetGroupID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetailByID", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetailByID", clientID, nil, "Asset group not found")
	}

	// Check if the user is a member of the asset group
	member, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroupID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if member == nil {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	// Get asset group details
	return assetGroup, nil
}

func (s *assetGroupService) DeleteAssetGroup(assetGroupID uint, clientID string) error {
	return nil
}

func (s *assetGroupService) AddMemberAssetGroup(req *request.AssetGroupMemberRequest, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	// Check if the user member permission is an admin or manager
	memberPermission, err := s.memberPermissionRepository.GetAdminOrManagePermissionsByUserID(user.UserID)

	if memberPermission == nil {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, nil, "User does not have permission to add members")
	}

	var hasPermission bool
	hasPermission = false

	for _, permission := range memberPermission {
		if permission.PermissionName == "Admin" || permission.PermissionName == "Manage" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, nil, "User does not have permission to add members")
	}

	// Check if the user is already a member of the asset group
	existingMember, _ := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(req.UserID, req.AssetGroupID)

	if existingMember != nil {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is already a member of this asset group")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(req.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupDetailByID", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logErrorWithNoReturn("GetAssetGroupDetailByID", clientID, nil, "Asset group not found")
	}

	// Check if the user exists
	member, err := s.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByID", clientID, err, "Failed to get user")
	}

	if member == nil {
		return logErrorWithNoReturn("GetUserByID", clientID, nil, "User not found")
	}

	groupMember := &assets.AssetGroupMember{
		UserID:       req.UserID,
		AssetGroupID: req.AssetGroupID,
		CreatedBy:    user.ClientID,
	}

	err = s.memberRepository.AddAssetGroupMember(groupMember, user.ClientID, member.ClientID)
	if err != nil {
		return logErrorWithNoReturn("AddAssetGroupMember", clientID, err, "Failed to add asset group member")
	}
	return nil
}
