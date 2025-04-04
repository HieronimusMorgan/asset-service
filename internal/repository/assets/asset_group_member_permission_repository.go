package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"gorm.io/gorm"
)

type AssetGroupMemberPermissionRepository interface {
	AddAssetGroupMemberPermission(asset *assets.AssetGroupMemberPermission) error
	RemoveAssetGroupMemberPermission(userID uint, assetGroupID uint, permissionID uint) error
	UpdateAssetGroupMemberPermission(asset *assets.AssetGroupMemberPermission) error
	GetAssetGroupMemberPermissionByID(assetGroupID uint) (*assets.AssetGroupMemberPermission, error)
	DeleteAssetGroupMemberPermission(assetGroupID uint) error
	GetAdminOrManagePermissionsByUserID(userID uint) ([]assets.AssetGroupPermission, error)
	GetAdminPermissionsByUserID(userID uint) ([]assets.AssetGroupPermission, error)
	GetAssetGroupMemberPermissionByUserIDAndGroupID(userID uint, assetGroupID uint) ([]assets.AssetGroupMemberPermission, error)
}

type assetGroupMemberPermissionRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
}

func NewAssetGroupMemberPermissionRepository(db gorm.DB, audit AssetAuditLogRepository) AssetGroupMemberPermissionRepository {
	return assetGroupMemberPermissionRepository{db: db, audit: audit}
}

func (r assetGroupMemberPermissionRepository) AddAssetGroupMemberPermission(asset *assets.AssetGroupMemberPermission) error {
	return r.db.Table(utils.TableAssetGroupMemberPermissionName).Create(asset).Error
}

func (r assetGroupMemberPermissionRepository) RemoveAssetGroupMemberPermission(userID uint, assetGroupID uint, permissionID uint) error {
	return r.db.Unscoped().Table(utils.TableAssetGroupMemberPermissionName).
		Where("user_id = ? AND asset_group_id = ? AND permission_id = ?", userID, assetGroupID, permissionID).
		Delete(&assets.AssetGroupMemberPermission{}).Error
}

func (r assetGroupMemberPermissionRepository) UpdateAssetGroupMemberPermission(asset *assets.AssetGroupMemberPermission) error {
	return r.db.Table(utils.TableAssetGroupMemberPermissionName).Save(asset).Error
}

func (r assetGroupMemberPermissionRepository) GetAssetGroupMemberPermissionByID(assetGroupID uint) (*assets.AssetGroupMemberPermission, error) {
	var asset assets.AssetGroupMemberPermission
	if err := r.db.Table(utils.TableAssetGroupMemberPermissionName).First(&asset, assetGroupID).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetGroupMemberPermissionRepository) DeleteAssetGroupMemberPermission(assetGroupID uint) error {
	return r.db.Table(utils.TableAssetGroupMemberPermissionName).Where("asset_group_id = ?", assetGroupID).Delete(&assets.AssetGroupMemberPermission{}).Error
}

func (r assetGroupMemberPermissionRepository) GetAdminOrManagePermissionsByUserID(userID uint) ([]assets.AssetGroupPermission, error) {
	var results []assets.AssetGroupPermission

	err := r.db.
		Table("asset_group_member_permission AS agmp").
		Select("agp.*").
		Joins("JOIN asset_group_permission AS agp ON agmp.permission_id = agp.permission_id").
		Where("agmp.user_id = ? AND (agp.permission_name = ? OR agp.permission_name = ?)", userID, "Admin", "Manage").
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r assetGroupMemberPermissionRepository) GetAdminPermissionsByUserID(userID uint) ([]assets.AssetGroupPermission, error) {
	var results []assets.AssetGroupPermission

	err := r.db.
		Table("asset_group_member_permission AS agmp").
		Select("agp.*").
		Joins("JOIN asset_group_permission AS agp ON agmp.permission_id = agp.permission_id").
		Where("agmp.user_id = ? AND (agp.permission_name = ? )", userID, "Admin").
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r assetGroupMemberPermissionRepository) GetAssetGroupMemberPermissionByUserIDAndGroupID(userID uint, assetGroupID uint) ([]assets.AssetGroupMemberPermission, error) {
	var assetGroupMemberPermission []assets.AssetGroupMemberPermission
	err := r.db.Table(utils.TableAssetGroupMemberPermissionName).
		Where("user_id = ? AND asset_group_id = ?", userID, assetGroupID).
		Find(&assetGroupMemberPermission).Error
	if err != nil {
		return nil, err
	}
	return assetGroupMemberPermission, nil
}
