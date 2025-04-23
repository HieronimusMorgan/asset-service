package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"fmt"
	"gorm.io/gorm"
)

type AssetGroupMemberRepository interface {
	AddAssetGroupMember(asset *assets.AssetGroupMember, userClientID string, memberClientID string) error
	UpdateAssetGroupMember(asset *assets.AssetGroupMember) error
	GetAssetGroupMemberByID(assetGroupID uint) (*[]response.AssetGroupMemberResponse, error)
	RemoveAssetGroupMember(assetGroupID, userID uint) error
	GetAssetGroupMemberByUserIDAndGroupID(userID uint, groupID uint) (assets.AssetGroupMember, error)
	GetAssetGroupMemberByUserID(userID uint) (*assets.AssetGroupMember, error)
}

type assetGroupMemberRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
}

func NewAssetGroupMemberRepository(db gorm.DB, audit AssetAuditLogRepository) AssetGroupMemberRepository {
	return assetGroupMemberRepository{db: db, audit: audit}
}

func (r assetGroupMemberRepository) AddAssetGroupMember(member *assets.AssetGroupMember, userClientID string, memberClientID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var permission []assets.AssetGroupPermission

		err := tx.Table(utils.TableAssetGroupPermissionName).Where("permission_name ='Read-Write' OR permission_name = 'Read'").Find(&permission).Error
		if err != nil {
			return err
		}

		for _, p := range permission {
			permissionRecord := &assets.AssetGroupMemberPermission{
				AssetGroupID: member.AssetGroupID,
				UserID:       member.UserID,
				CreatedBy:    userClientID,
				PermissionID: p.PermissionID,
			}
			if err := tx.Table(utils.TableAssetGroupMemberPermissionName).Create(permissionRecord).Error; err != nil {
				return err
			}
		}

		if err := tx.Table(utils.TableAssetGroupMemberName).Create(member).Error; err != nil {
			return err
		}

		//get assets
		var asset []assets.Asset
		err = tx.Table(utils.TableAssetName).Where("user_client_id = ?", memberClientID).Find(&asset).Error
		if err != nil {
			return err
		}

		for _, asset := range asset {
			if err := tx.Table(utils.TableAssetGroupAssetName).Create(&assets.AssetGroupAsset{
				AssetGroupID: member.AssetGroupID,
				AssetID:      asset.AssetID,
				UserID:       member.UserID,
				CreatedBy:    userClientID,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r assetGroupMemberRepository) UpdateAssetGroupMember(asset *assets.AssetGroupMember) error {
	return r.db.Table(utils.TableAssetGroupMemberName).Save(asset).Error
}

func (r assetGroupMemberRepository) GetAssetGroupMemberByID(groupID uint) (*[]response.AssetGroupMemberResponse, error) {

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
		ORDER BY u.user_id ASC
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
		ORDER BY agm.user_id ASC
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

	return &members, nil
}

func (r assetGroupMemberRepository) RemoveAssetGroupMember(assetGroupID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Table(utils.TableAssetGroupMemberPermissionName).Where("asset_group_id = ? AND user_id = ?", assetGroupID, userID).Delete(&assets.AssetGroupMemberPermission{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Table(utils.TableAssetGroupMemberName).Where("asset_group_id = ? AND user_id = ?", assetGroupID, userID).Delete(&assets.AssetGroupMember{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Table(utils.TableAssetGroupAssetName).Where("asset_group_id = ? AND user_id = ?", assetGroupID, userID).Delete(&assets.AssetGroupAsset{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r assetGroupMemberRepository) GetAssetGroupMemberByUserIDAndGroupID(userID uint, groupID uint) (assets.AssetGroupMember, error) {
	var assetGroupMember assets.AssetGroupMember
	err := r.db.Table(utils.TableAssetGroupMemberName).Where("user_id = ? AND asset_group_id = ?", userID, groupID).First(&assetGroupMember).Error
	if err != nil {
		return assets.AssetGroupMember{}, err
	}
	return assetGroupMember, nil
}

func (r assetGroupMemberRepository) GetAssetGroupMemberByUserID(userID uint) (*assets.AssetGroupMember, error) {
	var assetGroupMember *assets.AssetGroupMember
	err := r.db.Table(utils.TableAssetGroupMemberName).Where("user_id = ?", userID).First(&assetGroupMember).Error
	if err != nil {
		return nil, err
	}
	return assetGroupMember, nil
}
