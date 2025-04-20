package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	repousers "asset-service/internal/repository/users"
	"asset-service/internal/utils"
	"errors"
)

type AssetGroupMemberService interface {
	AddAssetGroupMember(req *request.AssetGroupMemberRequest, clientID string) error
	RemoveMemberAssetGroup(memberRequest request.AssetGroupMemberRequest, clientID string) error
	GetListAssetGroupMember(assetGroupID uint, clientID string) (interface{}, error)
	LeaveMemberAssetGroup(assetGroupID uint, clientID string) error
}

type assetGroupMemberService struct {
	UserRepository                       repousers.UserRepository
	UserSettingRepository                repousers.UserSettingRepository
	AssetGroupMemberPermissionRepository repository.AssetGroupMemberPermissionRepository
	AssetGroupRepository                 repository.AssetGroupRepository
	AssetGroupMemberRepository           repository.AssetGroupMemberRepository
	AssetRepository                      repository.AssetRepository
	AssetGroupInvitation                 repository.AssetGroupInvitationRepository
	AssetAuditLogRepository              repository.AssetAuditLogRepository
	Redis                                utils.RedisService
}

func NewAssetGroupMemberService(
	userRepository repousers.UserRepository,
	UserSettingRepository repousers.UserSettingRepository,
	assetGroupRepository repository.AssetGroupRepository,
	AssetGroupMemberPermissionRepository repository.AssetGroupMemberPermissionRepository,
	assetCategoryRepository repository.AssetGroupMemberRepository,
	AssetRepository repository.AssetRepository,
	AssetGroupInvitation repository.AssetGroupInvitationRepository,
	AssetAuditLogRepository repository.AssetAuditLogRepository,
	redis utils.RedisService) AssetGroupMemberService {
	return &assetGroupMemberService{
		UserRepository:                       userRepository,
		UserSettingRepository:                UserSettingRepository,
		AssetGroupMemberPermissionRepository: AssetGroupMemberPermissionRepository,
		AssetGroupRepository:                 assetGroupRepository,
		AssetGroupMemberRepository:           assetCategoryRepository,
		AssetRepository:                      AssetRepository,
		AssetGroupInvitation:                 AssetGroupInvitation,
		AssetAuditLogRepository:              AssetAuditLogRepository,
		Redis:                                redis,
	}
}

func (s *assetGroupMemberService) AddAssetGroupMember(req *request.AssetGroupMemberRequest, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	userPermission, err := s.AssetGroupMemberPermissionRepository.GetAdminOrManagePermissionsByUserID(user.UserID)

	if userPermission == nil {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, errors.New("user does not have permission to add members"), "User does not have permission to add members")
	}

	var hasPermission bool
	hasPermission = false

	for _, permission := range userPermission {
		if permission.PermissionName == "Admin" || permission.PermissionName == "Manage" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, errors.New("user does not have permission to add members"), "User does not have permission to add members")
	}

	existingUser, err := s.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByID", clientID, err, "Failed to get user")
	}

	if existingUser == nil {
		return logErrorWithNoReturn("GetUserByID", clientID, errors.New("user not found"), "User not found")
	}

	userSetting, err := s.UserSettingRepository.GetUserSettingByUserID(existingUser.UserID)
	if err != nil {
		return logErrorWithNoReturn("GetUserSettingByUserID", clientID, err, "Failed to get user setting")
	}

	if userSetting == nil {
		return logErrorWithNoReturn("GetUserSettingByUserID", clientID, errors.New("user setting not found"), "User setting not found")
	}

	if userSetting.GroupInviteType != 1 {
		existingMember, _ := s.AssetGroupMemberRepository.GetAssetGroupMemberByUserIDAndGroupID(req.UserID, req.AssetGroupID)

		if existingMember.AssetGroupID != 0 {
			return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, errors.New("user is already a member of this asset group"), "User is already a member of this asset group")
		}

		assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(req.AssetGroupID)
		if err != nil {
			return logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
		}

		if assetGroup == nil {
			return logErrorWithNoReturn("GetAssetGroupDetail", clientID, errors.New("asset group not found"), "Asset group not found")
		}

		member, err := s.UserRepository.GetUserByID(req.UserID)
		if err != nil {
			return logErrorWithNoReturn("GetUserByID", clientID, err, "Failed to get user")
		}

		if member == nil {
			return logErrorWithNoReturn("GetUserByID", clientID, errors.New("user not found"), "User not found")
		}

		inviteToken, err := utils.GenerateInviteToken()
		if err != nil {
			return logErrorWithNoReturn("GenerateInviteToken", clientID, err, "Failed to generate invite token")
		}
		invitation := &assets.AssetGroupInvitation{
			AssetGroupID:     assetGroup.AssetGroupID,
			InvitedUserID:    member.UserID,
			InvitedUserToken: inviteToken,
			InvitedByUserID:  user.UserID,
		}
		err = s.AssetGroupInvitation.AddAssetGroupInvitation(invitation)
		if err != nil {
			return logErrorWithNoReturn("AddAssetGroupInvitation", clientID, err, "Failed to add asset group invitation")
		}
	} else {
		existingMember, _ := s.AssetGroupMemberRepository.GetAssetGroupMemberByUserIDAndGroupID(req.UserID, req.AssetGroupID)

		if existingMember.AssetGroupID != 0 {
			return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, errors.New("user not found"), "User is already a member of this asset group")
		}

		assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(req.AssetGroupID)
		if err != nil {
			return logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
		}

		if assetGroup == nil {
			return logErrorWithNoReturn("GetAssetGroupDetail", clientID, errors.New("asset group not found"), "Asset group not found")
		}

		member, err := s.UserRepository.GetUserByID(req.UserID)
		if err != nil {
			return logErrorWithNoReturn("GetUserByID", clientID, err, "Failed to get user")
		}

		if member == nil {
			return logErrorWithNoReturn("GetUserByID", clientID, errors.New("user not found"), "User not found")
		}

		groupMember := &assets.AssetGroupMember{
			UserID:       req.UserID,
			AssetGroupID: req.AssetGroupID,
			CreatedBy:    user.ClientID,
		}

		err = s.AssetGroupMemberRepository.AddAssetGroupMember(groupMember, user.ClientID, member.ClientID)
		if err != nil {
			return logErrorWithNoReturn("AddAssetGroupMember", clientID, err, "Failed to add asset group member")
		}
	}
	return nil
}

func (s *assetGroupMemberService) RemoveMemberAssetGroup(memberRequest request.AssetGroupMemberRequest, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	userPermission, err := s.AssetGroupMemberPermissionRepository.GetAdminPermissionsByUserID(user.UserID)

	if userPermission == nil {
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, errors.New("user does not have permission to add members"), "User does not have permission to add members")
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
		return logErrorWithNoReturn("GetAssetGroupMemberPermissionByUserID", clientID, errors.New("user does not have permission to add members"), "User does not have permission to add members")
	}

	// Check if the asset group exists
	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(memberRequest.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logErrorWithNoReturn("GetAssetGroupDetail", clientID, errors.New("asset group not found"), "Asset group not found")
	}

	member, err := s.UserRepository.GetUserByID(memberRequest.UserID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByID", clientID, err, "Failed to get user")
	}

	if member == nil {
		return logErrorWithNoReturn("GetUserByID", clientID, errors.New("user not found"), "User not found")
	}

	// Check if the user is a member of the asset group
	existingMember, err := s.AssetGroupMemberRepository.GetAssetGroupMemberByUserIDAndGroupID(memberRequest.UserID, memberRequest.AssetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if existingMember.AssetGroupID == 0 {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, errors.New("user is not a member of this asset group"), "User is not a member of this asset group")
	}

	err = s.AssetGroupMemberRepository.RemoveAssetGroupMember(memberRequest.AssetGroupID, memberRequest.UserID)
	if err != nil {
		return logErrorWithNoReturn("RemoveAssetGroupMember", clientID, err, "Failed to remove asset group member")
	}

	return nil
}

func (s *assetGroupMemberService) LeaveMemberAssetGroup(assetGroupID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	existingMember, err := s.AssetGroupMemberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroupID)
	if err != nil {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if existingMember.AssetGroupID == 0 {
		return logErrorWithNoReturn("GetAssetGroupMemberByUserIDAndGroupID", clientID, errors.New("user is not a member of this asset group"), "User is not a member of this asset group")
	}

	err = s.AssetGroupMemberRepository.RemoveAssetGroupMember(assetGroupID, user.UserID)
	if err != nil {
		return logErrorWithNoReturn("RemoveAssetGroupMember", clientID, err, "Failed to remove asset group member")
	}
	return nil
}

func (s *assetGroupMemberService) GetListAssetGroupMember(assetGroupID uint, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logError("GetUserByClientID", clientID, err, "Failed to get user data")
	}

	assetGroup, err := s.AssetGroupRepository.GetAssetGroupByID(assetGroupID)
	if err != nil {
		return logError("GetAssetGroupDetail", clientID, err, "Failed to get asset group")
	}

	if assetGroup == nil {
		return logError("GetAssetGroupDetail", clientID, errors.New("asset group not found"), "Asset group not found")
	}

	existingMember, err := s.AssetGroupMemberRepository.GetAssetGroupMemberByUserIDAndGroupID(user.UserID, assetGroupID)
	if err != nil {
		return logError("GetAssetGroupMemberByUserIDAndGroupID", clientID, err, "Failed to get asset group member")
	}

	if existingMember.AssetGroupID == 0 {
		return logError("GetAssetGroupMemberByUserIDAndGroupID", clientID, errors.New("user is not a member of this asset group"), "User is not a member of this asset group")
	}

	members, err := s.AssetGroupMemberRepository.GetAssetGroupMemberByID(assetGroupID)
	if err != nil {
		return logError("GetAssetGroupMemberByID", clientID, err, "Failed to get asset group members")
	}

	return members, nil
}
