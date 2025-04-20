package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	users "asset-service/internal/repository/users"
	"asset-service/internal/utils"
	"errors"
	"github.com/rs/zerolog/log"
)

type AssetGroupService interface {
	AddAssetGroup(assetRequest *request.AssetGroupRequest, clientID string) (interface{}, error)
	AddInvitationAssetGroup(assetGroupID uint, clientID string) (interface{}, error)
	RemoveInvitationAssetGroup(assetGroupID uint, clientID string) error
	UpdateAssetGroup(assetGroupID uint, req *request.AssetGroupRequest, clientID string) (interface{}, error)
	GetAssetGroupDetail(clientID string) (interface{}, error)
	GetAssetGroupAssetByAssetGroupID(assetGroupID uint, clientID string) (interface{}, error)
	DeleteAssetGroup(assetGroupID uint, clientID string) error
	InviteMemberAssetGroup(req *request.AssetGroupMemberRequest, clientID string) error
	RemoveMemberAssetGroup(memberRequest request.AssetGroupMemberRequest, clientID string) error
	AddPermissionMemberAssetGroup(req *request.ChangeAssetGroupPermissionRequest, clientID string) error
	RemovePermissionMemberAssetGroup(req *request.ChangeAssetGroupPermissionRequest, clientID string) error
	GetListAssetGroupAsset(assetGroupID uint, clientID string) ([]response.AssetResponse, error)
	UpdateStockAssetGroupAsset(isAdded bool, req request.ChangeAssetStockRequest, clientID string) (interface{}, error)
}

type assetGroupService struct {
	UserRepository             users.UserRepository
	AssetGroupRepository       repository.AssetGroupRepository
	permissionRepository       repository.AssetGroupPermissionRepository
	memberPermissionRepository repository.AssetGroupMemberPermissionRepository
	memberRepository           repository.AssetGroupMemberRepository
	assetGroupAssetRepository  repository.AssetGroupAssetRepository
	AssetRepository            repository.AssetRepository
	AssetStockRepository       repository.AssetStockRepository
	AssetAuditLogRepository    repository.AssetAuditLogRepository
	Redis                      utils.RedisService
}

func NewAssetGroupService(UserRepository users.UserRepository, AssetGroupRepository repository.AssetGroupRepository, permissionRepository repository.AssetGroupPermissionRepository, memberPermissionRepository repository.AssetGroupMemberPermissionRepository, memberRepository repository.AssetGroupMemberRepository, assetGroupAssetRepository repository.AssetGroupAssetRepository, AssetRepository repository.AssetRepository, AssetStockRepository repository.AssetStockRepository, AssetAuditLogRepository repository.AssetAuditLogRepository, redis utils.RedisService) AssetGroupService {
	return &assetGroupService{
		UserRepository:             UserRepository,
		AssetGroupRepository:       AssetGroupRepository,
		permissionRepository:       permissionRepository,
		memberPermissionRepository: memberPermissionRepository,
		memberRepository:           memberRepository,
		assetGroupAssetRepository:  assetGroupAssetRepository,
		AssetRepository:            AssetRepository,
		AssetStockRepository:       AssetStockRepository,
		AssetAuditLogRepository:    AssetAuditLogRepository,
		Redis:                      redis,
	}
}

func (s *assetGroupService) AddAssetGroup(assetRequest *request.AssetGroupRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logError("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	if user.UserID != assetRequest.OwnerUserID {
		return logError("GetUserByClientID", clientID, nil, "User ID does not match the owner user ID")
	}

	//check if the user is already an owner or member of the asset group
	group, _ := s.AssetGroupRepository.GetAssetGroupByOwnerUserID(assetRequest.OwnerUserID)
	if len(group) != 0 {
		return logError("GetAssetGroupByOwnerUserID", clientID, nil, "User is already an owner or member of this asset group")
	}

	assetGroupMember, _ := s.memberRepository.GetAssetGroupMemberByUserID(user.UserID)
	if assetGroupMember != nil {
		return logError("GetAssetGroupMemberByUserID", clientID, nil, "User is already a member of this asset group")
	}

	// Check if the asset group name is empty
	if assetRequest.AssetGroupName == "" {
		return logError("AddAssetGroup", clientID, nil, "Asset group name cannot be empty")
	}

	// Add asset group
	assetGroup := &assets.AssetGroup{
		AssetGroupName: assetRequest.AssetGroupName,
		Description:    assetRequest.Description,
		OwnerUserID:    user.UserID,
		CreatedBy:      user.ClientID,
		UpdatedBy:      user.ClientID,
	}

	err = s.AssetGroupRepository.AddAssetGroup(assetGroup, clientID, user)
	if err != nil {
		return logError("AddAssetGroup", clientID, err, "Failed to add asset group")
	}
	return struct {
		AssetGroupID   uint   `json:"asset_group_id"`
		AssetGroupName string `json:"asset_group_name"`
		Description    string `json:"description"`
	}{
		AssetGroupID:   assetGroup.AssetGroupID,
		AssetGroupName: assetGroup.AssetGroupName,
		Description:    assetGroup.Description,
	}, nil
}

func (s *assetGroupService) AddInvitationAssetGroup(assetGroupID uint, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logError("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(assetGroupID)
	if err != nil {
		return logError("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logError("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	// Check if the user is a member of the asset group
	member, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroupID)
	if err != nil {
		return logError("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}
	if member.AssetGroupID == 0 {
		return logError("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	// Check if the user has permission to add invitation token
	permissions, err := s.memberPermissionRepository.GetAdminOrManagePermissionsByUserID(user.UserID)
	if err != nil {
		return logError("GetAssetGroupPermissionByUserID", clientID, err, "Failed to get asset group permission")
	}

	if len(permissions) == 0 {
		return logError("GetAssetGroupPermissionByUserID", clientID, nil, "User does not have permission to add invitation token")
	}

	var hasPermission bool
	hasPermission = false
	for _, permission := range permissions {
		if permission.PermissionName == "Admin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return logError("GetAssetGroupPermissionByUserID", clientID, nil, "User does not have permission to add invitation token")
	}

	invitationToken, err := utils.GenerateInviteToken()
	if err != nil {
		return logError("GenerateInviteToken", clientID, err, "Failed to generate invitation token")
	}
	err = s.AssetGroupRepository.AddInvitationToken(assetGroupID, invitationToken, user.ClientID)
	if err != nil {
		return logError("AddInvitationToken", clientID, err, "Failed to add invitation token")
	}

	return struct {
		InvitationToken string `json:"invitation_token"`
	}{
		InvitationToken: invitationToken,
	}, nil
}

func (s *assetGroupService) RemoveInvitationAssetGroup(assetGroupID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(assetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	// Check if the user is a member of the asset group
	member, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}
	if member.AssetGroupID == 0 {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	// Check if the user has permission to add invitation token
	permissions, err := s.memberPermissionRepository.GetAdminOrManagePermissionsByUserID(user.UserID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupPermissionByUserID", clientID, err, "Failed to get asset group permission")
	}

	if len(permissions) == 0 {
		return logErrorWithNoReturn("GetAssetGroupPermissionByUserID", clientID, nil, "User does not have permission to add invitation token")
	}

	var hasPermission bool
	hasPermission = false
	for _, permission := range permissions {
		if permission.PermissionName == "Admin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return logErrorWithNoReturn("GetAssetGroupPermissionByUserID", clientID, nil, "User does not have permission to add invitation token")
	}

	err = s.AssetGroupRepository.RemoveInvitationToken(assetGroupID, user.ClientID)
	if err != nil {
		return logErrorWithNoReturn("AddInvitationToken", clientID, err, "Failed to add invitation token")
	}

	return nil
}

func (s *assetGroupService) UpdateAssetGroup(assetGroupID uint, req *request.AssetGroupRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}
	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	member, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroupID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if member.AssetGroupID == 0 {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	permissions, err := s.memberPermissionRepository.GetAdminOrManagePermissionsByUserID(user.UserID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupPermissionByUserID", clientID, err, "Failed to get asset group permission")
	}
	if len(permissions) == 0 {
		return nil, logErrorWithNoReturn("GetAssetGroupPermissionByUserID", clientID, nil, "User does not have permission to update asset group")
	}

	var hasPermission bool
	hasPermission = false
	for _, permission := range permissions {
		if permission.PermissionName == "Admin" || permission.PermissionName == "Manage" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, logErrorWithNoReturn("GetAssetGroupPermissionByUserID", clientID, nil, "User does not have permission to update asset group")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(assetGroupID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	// Update asset group
	if assetGroup.AssetGroupName == "" {
		return nil, logErrorWithNoReturn("UpdateAssetGroup", clientID, nil, "Asset group name cannot be empty")
	}
	assetGroup.AssetGroupName = req.AssetGroupName
	assetGroup.Description = req.Description
	assetGroup.UpdatedBy = user.ClientID
	err = s.AssetGroupRepository.UpdateAssetGroup(assetGroup)
	if err != nil {
		return nil, logErrorWithNoReturn("UpdateAssetGroup", clientID, err, "Failed to update asset group")
	}

	return struct {
		AssetGroupID   uint   `json:"asset_group_id"`
		AssetGroupName string `json:"asset_group_name"`
		Description    string `json:"description"`
	}{
		AssetGroupID:   assetGroup.AssetGroupID,
		AssetGroupName: assetGroup.AssetGroupName,
		Description:    assetGroup.Description,
	}, nil
}

func (s *assetGroupService) GetAssetGroupDetail(clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupDetailByUserID(user.UserID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	// Check if the user is a member of the asset group
	member, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroup.AssetGroupID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if member.AssetGroupID == 0 {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	// Get asset group details
	return assetGroup, nil
}

func (s *assetGroupService) GetAssetGroupAssetByAssetGroupID(assetGroupID uint, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupDetailByUserID(user.UserID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	// Check if the user is a member of the asset group
	member, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroupID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if member.AssetGroupID == 0 {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	// Get asset group assets
	assetsList, err := s.assetGroupAssetRepository.GetAssetGroupAssetByID(assetGroup.AssetGroupID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetsByAssetGroup", clientID, err, "Failed to get assets by asset group")
	}

	return assetsList, nil
}

func (s *assetGroupService) DeleteAssetGroup(assetGroupID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}
	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}
	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(assetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	// Check if the user is a member of the asset group
	member, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if member.AssetGroupID == 0 {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	// Check if the user has permission to delete the asset group
	permissions, err := s.memberPermissionRepository.GetAdminOrManagePermissionsByUserID(user.UserID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupPermissionByUserID", clientID, err, "Failed to get asset group permission")
	}

	if len(permissions) == 0 {
		return logErrorWithNoReturn("GetAssetGroupPermissionByUserID", clientID, nil, "User does not have permission to delete asset group")
	}

	var hasPermission bool
	hasPermission = false
	for _, permission := range permissions {
		if permission.PermissionName == "Admin" || permission.PermissionName == "Manage" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return logErrorWithNoReturn("GetAssetGroupPermissionByUserID", clientID, nil, "User does not have permission to delete asset group")
	}

	// Delete asset group
	err = s.AssetGroupRepository.DeleteAssetGroup(assetGroupID, user.UserID)
	if err != nil {
		return logErrorWithNoReturn("DeleteAssetGroup", clientID, err, "Failed to delete asset group")
	}

	return nil
}

func (s *assetGroupService) InviteMemberAssetGroup(req *request.AssetGroupMemberRequest, clientID string) error {
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

	if existingMember.AssetGroupID != 0 {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is already a member of this asset group")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(req.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
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

func (s *assetGroupService) RemoveMemberAssetGroup(memberRequest request.AssetGroupMemberRequest, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	userPermission, err := s.memberPermissionRepository.GetAdminPermissionsByUserID(user.UserID)

	if userPermission == nil {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, nil, "User does not have permission to add members")
	}

	var hasPermission bool
	hasPermission = false

	for _, permission := range userPermission {
		if permission.PermissionName == "Admin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, nil, "User does not have permission to add members")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(memberRequest.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	member, err := s.UserRepository.GetUserByID(memberRequest.UserID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByID", clientID, err, "Failed to get user")
	}

	if member == nil {
		return logErrorWithNoReturn("GetUserByID", clientID, nil, "User not found")
	}

	// Check if the user is a member of the asset group
	existingMember, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(memberRequest.UserID, memberRequest.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if existingMember.AssetGroupID == 0 {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	err = s.memberRepository.RemoveAssetGroupMember(memberRequest.AssetGroupID, memberRequest.UserID)
	if err != nil {
		return logErrorWithNoReturn("RemoveAssetGroupMember", clientID, err, "Failed to remove asset group member")
	}
	return nil
}

func (s *assetGroupService) AddPermissionMemberAssetGroup(req *request.ChangeAssetGroupPermissionRequest, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	userPermission, err := s.memberPermissionRepository.GetAdminPermissionsByUserID(user.UserID)

	if userPermission == nil {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, nil, "User does not have permission to add members")
	}

	var hasPermission bool
	hasPermission = false

	for _, permission := range userPermission {
		if permission.PermissionName == "Admin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, nil, "User does not have permission to add members")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(req.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	// Check if the user exists
	member, err := s.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByID", clientID, err, "Failed to get user")
	}

	if member == nil {
		return logErrorWithNoReturn("GetUserByID", clientID, nil, "User not found")
	}

	permission, err := s.permissionRepository.GetAssetGroupPermissionByID(req.PermissionID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupPermissionByID", clientID, err, "Failed to get asset group permission")
	}

	if permission == nil {
		return logErrorWithNoReturn("GetAssetGroupPermissionByID", clientID, nil, "Asset group permission not found")
	}

	// Check if the user already has the permission
	existingPermission, err := s.memberPermissionRepository.GetAssetGroupMemberPermissionByUserIDAndGroupID(req.UserID, req.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserIDAndGroupID", clientID, err, "Failed to get asset group member permission")
	}

	for _, permission := range existingPermission {
		if permission.PermissionID == req.PermissionID {
			return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserIDAndGroupID", clientID, nil, "User already has this permission")
		}
	}

	// Add permission to the user
	groupMemberPermission := &assets.AssetGroupMemberPermission{
		UserID:       req.UserID,
		AssetGroupID: req.AssetGroupID,
		PermissionID: req.PermissionID,
		CreatedBy:    user.ClientID,
	}

	err = s.memberPermissionRepository.AddAssetGroupMemberPermission(groupMemberPermission)
	if err != nil {
		return logErrorWithNoReturn("AddAssetGroupMemberPermission", clientID, err, "Failed to add asset group member permission")
	}

	return nil
}

func (s *assetGroupService) RemovePermissionMemberAssetGroup(req *request.ChangeAssetGroupPermissionRequest, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	userPermission, err := s.memberPermissionRepository.GetAdminPermissionsByUserID(user.UserID)

	if userPermission == nil {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, nil, "User does not have permission to add members")
	}

	var hasPermission bool
	hasPermission = false

	for _, permission := range userPermission {
		if permission.PermissionName == "Admin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, nil, "User does not have permission to add members")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(req.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	// Check if the user exists
	member, err := s.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByID", clientID, err, "Failed to get user")
	}

	if member == nil {
		return logErrorWithNoReturn("GetUserByID", clientID, nil, "User not found")
	}

	permission, err := s.permissionRepository.GetAssetGroupPermissionByID(req.PermissionID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupPermissionByID", clientID, err, "Failed to get asset group permission")
	}

	if permission == nil {
		return logErrorWithNoReturn("GetAssetGroupPermissionByID", clientID, nil, "Asset group permission not found")
	}

	// Check if the user already has the permission
	existingPermission, err := s.memberPermissionRepository.GetAssetGroupMemberPermissionByUserIDAndGroupID(req.UserID, req.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserIDAndGroupID", clientID, err, "Failed to get asset group member permission")
	}

	hasPermission = false
	for _, permission := range existingPermission {
		if permission.PermissionID == req.PermissionID {
			hasPermission = true
		}
	}

	if !hasPermission {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserIDAndGroupID", clientID, nil, "User does not have this permission")
	}

	err = s.memberPermissionRepository.RemoveAssetGroupMemberPermission(req.UserID, req.AssetGroupID, req.PermissionID)
	if err != nil {
		return logErrorWithNoReturn("AddAssetGroupMemberPermission", clientID, err, "Failed to add asset group member permission")
	}

	return nil
}

func (s *assetGroupService) GetListAssetGroupAsset(assetGroupID uint, clientID string) ([]response.AssetResponse, error) {
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
		return nil, logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return nil, logErrorWithNoReturn("GetAssetGroupDetail", clientID, nil, "Asset group not found")
	}

	// Check if the user is a member of the asset group
	member, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroupID)
	if err != nil {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if member.AssetGroupID == 0 {
		return nil, logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	asset, err := s.AssetRepository.GetListAssetsByAssetGroup(user.ClientID, assetGroupID)

	return asset, nil
}

func (s *assetGroupService) UpdateStockAssetGroupAsset(isAdded bool, req request.ChangeAssetStockRequest, clientID string) (interface{}, error) {
	// Step 1: Fetch user data from Redis
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user from Redis")
	}

	// Step 1: Fetch user data from database
	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logError("GetUserByClientID", clientID, err, "Failed to get user by client ID")
	}

	// Step 1: Check if user is a member of the asset group
	member, err := s.memberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, req.AssetGroupID)
	if err != nil {
		return logError("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if member.AssetGroupID == 0 {
		return logError("GetAssetGroupMemberByUserIDAndGroupID", clientID, nil, "User is not a member of this asset group")
	}

	// Step 2: Retrieve asset and stock data
	asset, err := s.AssetRepository.GetAssetByAssetGroupID(req.AssetID, req.AssetGroupID)
	if err != nil {
		return logError("GetAsset", clientID, err, "Failed to get asset by ID")
	}

	oldAssetStock, err := s.AssetStockRepository.GetAssetStockByAssetIDAndAssetGroupID(asset.AssetID, req.AssetGroupID)
	if err != nil {
		return logError("GetAssetStockByAssetID", clientID, err, "Failed to get asset stock by asset ID")
	}

	log.Info().
		Uint("assetID", req.AssetID).
		Int("Previous Stock", oldAssetStock.LatestQuantity).
		Msg("Retrieved current stock")

	// Step 3: Determine new stock quantity
	var stockType string
	var latestQuantity int

	if isAdded {
		stockType = "INCREASE"
		latestQuantity = oldAssetStock.LatestQuantity + req.Stock
	} else {
		stockType = "DECREASE"
		if oldAssetStock.LatestQuantity < req.Stock {
			return logError("UpdateStockAsset", clientID, errors.New("insufficient stock"), "Stock cannot be negative")
		}
		latestQuantity = oldAssetStock.LatestQuantity - req.Stock
	}

	// Step 4: Create stock update struct
	newAssetStock := &assets.AssetStock{
		AssetID:         asset.AssetID,
		UserClientID:    data.ClientID,
		InitialQuantity: oldAssetStock.InitialQuantity,
		LatestQuantity:  latestQuantity,
		Quantity:        req.Stock,
		ChangeType:      stockType,
		Reason:          req.Reason,
		UpdatedBy:       data.ClientID,
	}

	// Step 5: Update stock in a transaction
	err = s.AssetStockRepository.UpdateAssetStockByAssetGroupID(newAssetStock, req.AssetGroupID, clientID)
	if err != nil {
		return logError("UpdateAssetStock", clientID, err, "Failed to update asset stock")
	}

	log.Info().
		Uint("assetID", req.AssetID).
		Int("Updated Stock", latestQuantity).
		Str("Change Type", stockType).
		Msg("Stock updated successfully")

	// Step 6: Log stock change in audit log
	err = s.AssetAuditLogRepository.AfterUpdateAssetStock(*oldAssetStock, newAssetStock)
	if err != nil {
		return logError("AfterUpdateAssetStock", clientID, err, "Failed to update asset stock log")
	}

	// Step 7: Return updated stock response
	return response.AssetStockResponse{
		StockID:         newAssetStock.StockID,
		AssetID:         newAssetStock.AssetID,
		InitialQuantity: newAssetStock.InitialQuantity,
		LatestQuantity:  newAssetStock.LatestQuantity,
		ChangeType:      newAssetStock.ChangeType,
		Quantity:        newAssetStock.Quantity,
	}, nil
}
