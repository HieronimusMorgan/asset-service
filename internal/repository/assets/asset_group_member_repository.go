package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"gorm.io/gorm"
)

type AssetGroupMemberRepository interface {
	AddAssetGroupMember(asset *assets.AssetGroupMember, userClientID string, memberClientID string) error
	UpdateAssetGroupMember(asset *assets.AssetGroupMember) error
	GetAssetGroupMemberByID(assetGroupID uint) (*assets.AssetGroupMember, error)
	DeleteAssetGroupMember(assetGroupID, userID uint) error
	GetAssetGroupMemberByUserIDAndGroupID(userID uint, groupID uint) (interface{}, error)
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

func (r assetGroupMemberRepository) GetAssetGroupMemberByID(assetGroupID uint) (*assets.AssetGroupMember, error) {
	var asset *assets.AssetGroupMember
	err := r.db.Table(utils.TableAssetGroupMemberName).Where("asset_group_id = ?", assetGroupID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (r assetGroupMemberRepository) DeleteAssetGroupMember(assetGroupID, userID uint) error {
	return r.db.Table(utils.TableAssetGroupMemberName).Where("asset_group_id = ? AND user_id = ?", assetGroupID, userID).Delete(&assets.AssetGroupMember{}).Error
}

func (r assetGroupMemberRepository) GetAssetGroupMemberByUserIDAndGroupID(userID uint, groupID uint) (interface{}, error) {
	var assetGroupMember assets.AssetGroupMember
	err := r.db.Table(utils.TableAssetGroupMemberName).Where("user_id = ? AND asset_group_id = ?", userID, groupID).First(&assetGroupMember).Error
	if err != nil {
		return nil, err
	}
	return assetGroupMember, nil
}
