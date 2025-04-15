package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	"asset-service/internal/utils"
	"fmt"
	"gorm.io/gorm"
)

type AssetGroupRepository interface {
	AddAssetGroup(assetGroup *assets.AssetGroup, clientID string, user *user.Users) error
	AddInvitationToken(assetGroupID uint, token string, clientID string) error
	RemoveInvitationToken(assetGroupID uint, clientID string) error
	UpdateCurrentUsesInvitationToken(assetGroupID uint, clientID string) error
	UpdateAssetGroup(asset *assets.AssetGroup, userID uint) error
	GetAssetGroupByID(assetGroupID uint) (*assets.AssetGroup, error)
	GetAssetGroupDetailByID(groupID uint) (*response.AssetGroupDetailResponse, error)
	GetAssetGroupByOwnerUserID(id uint) ([]assets.AssetGroup, error)
	DeleteAssetGroup(assetGroupID uint, userID uint) error
	GetAssetGroupByInvitationToken(invitationToken string) (*assets.AssetGroup, error)
}

type assetGroupRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
}

func NewAssetGroupRepository(db gorm.DB, audit AssetAuditLogRepository) AssetGroupRepository {
	return assetGroupRepository{db: db, audit: audit}
}

func (r assetGroupRepository) AddAssetGroup(assetGroup *assets.AssetGroup, clientID string, user *user.Users) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableAssetGroupName).Create(&assetGroup).Error; err != nil {
			return err
		}

		var permission []assets.AssetGroupPermission

		err := tx.Table(utils.TableAssetGroupPermissionName).Find(&permission).Error
		if err != nil {
			return err
		}

		for _, p := range permission {
			permissionRecord := &assets.AssetGroupMemberPermission{
				AssetGroupID: assetGroup.AssetGroupID,
				UserID:       user.UserID,
				CreatedBy:    user.ClientID,
				PermissionID: p.PermissionID,
			}
			if err := tx.Table(utils.TableAssetGroupMemberPermissionName).Create(permissionRecord).Error; err != nil {
				return err
			}
		}

		groupMember := &assets.AssetGroupMember{
			AssetGroupID: assetGroup.AssetGroupID,
			UserID:       user.UserID,
			CreatedBy:    user.ClientID,
		}

		if err := tx.Table(utils.TableAssetGroupMemberName).Create(&groupMember).Error; err != nil {
			return err
		}

		//get assets
		var asset []assets.Asset
		err = tx.Table(utils.TableAssetName).Where("user_client_id = ?", clientID).Find(&asset).Error
		if err != nil {
			return err
		}

		for _, asset := range asset {
			if err := tx.Table(utils.TableAssetGroupAssetName).Create(&assets.AssetGroupAsset{
				AssetGroupID: assetGroup.AssetGroupID,
				AssetID:      asset.AssetID,
				UserID:       user.UserID,
				CreatedBy:    user.ClientID,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})

}

func (r assetGroupRepository) AddInvitationToken(assetGroupID uint, token string, clientID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var assetGroup assets.AssetGroup
		if err := tx.Table(utils.TableAssetGroupName).Where("asset_group_id = ?", assetGroupID).First(&assetGroup).Error; err != nil {
			return err
		}

		invitationToken := token
		maxUses := 10
		currentUses := 0
		assetGroup.InvitationToken = &invitationToken
		assetGroup.MaxUses = &maxUses
		assetGroup.CurrentUses = &currentUses
		assetGroup.UpdatedBy = clientID

		if err := tx.Table(utils.TableAssetGroupName).Save(&assetGroup).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r assetGroupRepository) RemoveInvitationToken(assetGroupID uint, clientID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var assetGroup assets.AssetGroup
		if err := tx.Table(utils.TableAssetGroupName).Where("asset_group_id = ?", assetGroupID).First(&assetGroup).Error; err != nil {
			return err
		}

		assetGroup.InvitationToken = nil
		assetGroup.MaxUses = nil
		assetGroup.CurrentUses = nil
		assetGroup.UpdatedBy = clientID

		if err := tx.Table(utils.TableAssetGroupName).Save(&assetGroup).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r assetGroupRepository) UpdateCurrentUsesInvitationToken(assetGroupID uint, clientID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var assetGroup assets.AssetGroup
		if err := tx.Table(utils.TableAssetGroupName).Where("asset_group_id = ?", assetGroupID).First(&assetGroup).Error; err != nil {
			return err
		}

		currentUses := *assetGroup.CurrentUses + 1
		assetGroup.CurrentUses = &currentUses
		assetGroup.UpdatedBy = clientID

		if err := tx.Table(utils.TableAssetGroupName).Save(&assetGroup).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r assetGroupRepository) UpdateAssetGroup(asset *assets.AssetGroup, userID uint) error {
	return r.db.Table(utils.TableAssetGroupName).Save(asset).Error
}

func (r assetGroupRepository) GetAssetGroupByID(assetGroupID uint) (*assets.AssetGroup, error) {
	var asset assets.AssetGroup
	if err := r.db.Table(utils.TableAssetGroupName).Where("asset_group_id = ?", assetGroupID).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetGroupRepository) GetAssetGroupDetailByID(groupID uint) (*response.AssetGroupDetailResponse, error) {

	var groupRow struct {
		AssetGroupID   uint
		AssetGroupName string
		Description    string
		OwnerUserID    uint
		OwnerName      string
	}

	query := `
	SELECT 
		ag.asset_group_id,
		ag.asset_group_name,
		ag.description,
		ag.owner_user_id AS owner_user_id,
		u.full_name AS owner_name
	FROM asset_group ag
	LEFT JOIN users u ON ag.owner_user_id = u.user_id
	WHERE ag.asset_group_id = ?
`

	if err := r.db.Raw(query, groupID).Scan(&groupRow).Error; err != nil {
		return nil, fmt.Errorf("failed to get asset group info: %w", err)
	}
	var group response.AssetGroupDetailResponse

	group.AssetGroupID = groupRow.AssetGroupID
	group.AssetGroupName = groupRow.AssetGroupName
	group.Description = groupRow.Description
	group.OwnerUserID = groupRow.OwnerUserID
	group.OwnerName = groupRow.OwnerName
	group.Member = []response.AssetGroupMemberResponse{}

	type memberRow struct {
		UserID         uint
		Username       string
		FullName       string
		ProfilePicture string
	}

	var memberRows []memberRow
	memberQuery := `
		SELECT 
			u.user_id,
			u.username,
			u.full_name,
			u.profile_picture
		FROM asset_group_member AS agm
		LEFT JOIN users AS u ON agm.user_id = u.user_id
		WHERE u.deleted_at IS NULL AND agm.asset_group_id = ?
	`
	if err := r.db.Raw(memberQuery, groupID).Scan(&memberRows).Error; err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}

	var members []response.AssetGroupMemberResponse
	for _, row := range memberRows {
		members = append(members, response.AssetGroupMemberResponse{
			UserID:         row.UserID,
			Username:       row.Username,
			FullName:       row.FullName,
			ProfilePicture: row.ProfilePicture,
			Permission:     []response.AssetGroupMemberWithPermissionResponse{},
		})
	}

	type tempPerm struct {
		UserID         uint
		PermissionID   *uint
		PermissionName *string
	}
	var allPermissions []tempPerm
	permQuery := `
		SELECT 
			agm.user_id,
			agp.permission_id,
			agp.permission_name
		FROM asset_group_member AS agm
		LEFT JOIN asset_group_member_permission AS agmp 
			ON agm.user_id = agmp.user_id AND agm.asset_group_id = agmp.asset_group_id
		LEFT JOIN asset_group_permission AS agp 
			ON agp.permission_id = agmp.permission_id
		WHERE agm.asset_group_id = ?
	`
	if err := r.db.Raw(permQuery, groupID).Scan(&allPermissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get member permissions: %w", err)
	}

	permissionMap := make(map[uint][]response.AssetGroupMemberWithPermissionResponse)
	for _, p := range allPermissions {
		permissionMap[p.UserID] = append(permissionMap[p.UserID], response.AssetGroupMemberWithPermissionResponse{
			PermissionID:   p.PermissionID,
			PermissionName: p.PermissionName,
		})
	}

	for i := range members {
		members[i].Permission = permissionMap[members[i].UserID]
	}

	group.Member = members
	return &group, nil
}

func (r assetGroupRepository) GetAssetGroupByOwnerUserID(id uint) ([]assets.AssetGroup, error) {
	var assetGroups []assets.AssetGroup
	if err := r.db.Table(utils.TableAssetGroupName).Where("owner_user_id = ?", id).Find(&assetGroups).Error; err != nil {
		return nil, err
	}
	return assetGroups, nil
}

func (r assetGroupRepository) DeleteAssetGroup(assetGroupID uint, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Table(utils.TableAssetGroupMemberPermissionName).Where("asset_group_id = ?", assetGroupID).Delete(&assets.AssetGroupMemberPermission{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Table(utils.TableAssetGroupMemberName).Where("asset_group_id = ?", assetGroupID).Delete(&assets.AssetGroupMember{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Table(utils.TableAssetGroupAssetName).Where("asset_group_id = ?", assetGroupID).Delete(&assets.AssetGroupAsset{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Table(utils.TableAssetGroupName).Where("asset_group_id = ? AND owner_user_id = ?", assetGroupID, userID).Delete(&assets.AssetGroup{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r assetGroupRepository) GetAssetGroupByInvitationToken(invitationToken string) (*assets.AssetGroup, error) {
	var assetGroup assets.AssetGroup
	if err := r.db.Table(utils.TableAssetGroupName).Where("invitation_token = ?", invitationToken).First(&assetGroup).Error; err != nil {
		return nil, err
	}
	return &assetGroup, nil
}
