package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"gorm.io/gorm"
)

type AssetGroupPermissionRepository interface {
	AddAssetGroupPermission(asset *assets.AssetGroupPermission) error
	UpdateAssetGroupPermission(asset *assets.AssetGroupPermission) error
	GetAssetGroupPermissionByID(permissionID uint) (*assets.AssetGroupPermission, error)
	GetAssetGroupPermissionByUserID(userID uint) ([]assets.AssetGroupPermission, error)
	GetListAssetGroupPermission() (*[]assets.AssetGroupPermission, error)
	DeleteAssetGroupPermission(asset *assets.AssetGroupPermission) error
}

type assetGroupPermissionRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
}

func NewAssetGroupPermissionRepository(db gorm.DB, audit AssetAuditLogRepository) AssetGroupPermissionRepository {
	return assetGroupPermissionRepository{db: db, audit: audit}
}

func (r assetGroupPermissionRepository) AddAssetGroupPermission(asset *assets.AssetGroupPermission) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableAssetGroupPermissionName).Create(&asset).Error; err != nil {
			return err
		}
		if err := r.audit.AfterCreateAssetGroupPermission(tx, asset); err != nil {
			return err
		}
		return nil
	})
}

func (r assetGroupPermissionRepository) UpdateAssetGroupPermission(asset *assets.AssetGroupPermission) error {
	return r.db.Table(utils.TableAssetGroupPermissionName).Save(asset).Error
}

func (r assetGroupPermissionRepository) GetAssetGroupPermissionByID(permissionID uint) (*assets.AssetGroupPermission, error) {
	var asset assets.AssetGroupPermission
	if err := r.db.Table(utils.TableAssetGroupPermissionName).First(&asset, permissionID).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetGroupPermissionRepository) GetAssetGroupPermissionByUserID(userID uint) ([]assets.AssetGroupPermission, error) {
	var permissions []assets.AssetGroupPermission
	if err := r.db.Table(utils.TableAssetGroupPermissionName).
		Where("user_id = ?", userID).
		Order("user_id ASC").
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r assetGroupPermissionRepository) GetListAssetGroupPermission() (*[]assets.AssetGroupPermission, error) {
	var permissions []assets.AssetGroupPermission
	if err := r.db.Table(utils.TableAssetGroupPermissionName).
		Order("permission_id ASC").
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return &permissions, nil
}

func (r assetGroupPermissionRepository) DeleteAssetGroupPermission(permission *assets.AssetGroupPermission) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var permissionMember *[]assets.AssetGroupMemberPermission
		if err := tx.Table(utils.TableAssetGroupMemberPermissionName).Where("permission_id = ?", permission.PermissionID).Delete(&permissionMember).Error; err != nil {
			return err
		}

		if err := tx.Table(utils.TableAssetGroupPermissionName).Where("permission_id = ?", permission.PermissionID).Delete(&permission).Error; err != nil {
			return err
		}

		if err := r.audit.AfterDeleteAssetGroupPermission(tx, permission); err != nil {
			return err
		}
		return nil
	})
}
