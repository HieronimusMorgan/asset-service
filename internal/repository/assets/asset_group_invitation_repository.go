package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/jwt"
	"gorm.io/gorm"
)

type AssetGroupInvitationRepository interface {
	AddAssetGroupInvitation(asset *assets.AssetGroupInvitation) error
	DeleteAssetGroupInvitationByID(invitationID uint) error
	GetAssetGroupInvitationByID(invitationID uint) (*assets.AssetGroupInvitation, error)
	GetAssetGroupInvitationByInvitedUserID(userID uint) ([]assets.AssetGroupInvitation, error)
	GetAssetGroupInvitationByInvitedByUserID(userID uint) ([]assets.AssetGroupInvitation, error)
	GetListAssetGroupInvitation() (*[]assets.AssetGroupInvitation, error)
	UpdateAssetGroupInvitationByUserID(status string, userID uint) error
	UpdateAssetGroupInvitationByInvitationTokenAndUserID(status string, invitationToken, userID uint) error
	DeleteAssetGroupInvitationExpired() error
}

type assetGroupInvitationRepository struct {
	db gorm.DB
}

func NewAssetGroupInvitationRepository(db gorm.DB) AssetGroupInvitationRepository {
	return assetGroupInvitationRepository{db: db}
}

func (r assetGroupInvitationRepository) AddAssetGroupInvitation(asset *assets.AssetGroupInvitation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableAssetGroupInvitationName).Create(&asset).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r assetGroupInvitationRepository) DeleteAssetGroupInvitationByID(invitationID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableAssetGroupInvitationName).Where("invitation_id = ?", invitationID).Delete(&assets.AssetGroupInvitation{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r assetGroupInvitationRepository) GetAssetGroupInvitationByID(invitationID uint) (*assets.AssetGroupInvitation, error) {
	var asset assets.AssetGroupInvitation
	if err := r.db.Table(utils.TableAssetGroupInvitationName).First(&asset, invitationID).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetGroupInvitationRepository) GetAssetGroupInvitationByInvitedUserID(userID uint) ([]assets.AssetGroupInvitation, error) {
	var permissions []assets.AssetGroupInvitation
	if err := r.db.Table(utils.TableAssetGroupInvitationName).
		Where("invited_user_id = ?", userID).
		Order("invited_user_id ASC").
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r assetGroupInvitationRepository) GetAssetGroupInvitationByInvitedByUserID(userID uint) ([]assets.AssetGroupInvitation, error) {
	var permissions []assets.AssetGroupInvitation
	if err := r.db.Table(utils.TableAssetGroupInvitationName).Where("invited_by_user_id = ?", userID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r assetGroupInvitationRepository) GetListAssetGroupInvitation() (*[]assets.AssetGroupInvitation, error) {
	var permissions []assets.AssetGroupInvitation
	if err := r.db.Table(utils.TableAssetGroupInvitationName).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return &permissions, nil
}

func (r assetGroupInvitationRepository) UpdateAssetGroupInvitationByUserID(status string, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableAssetGroupInvitationName).Where("invited_user_id = ?", userID).Updates(map[string]interface{}{
			"status": status,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r assetGroupInvitationRepository) UpdateAssetGroupInvitationByInvitationTokenAndUserID(status string, invitationToken, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableAssetGroupInvitationName).Where("invitation_token = ? AND invited_user_id = ?", invitationToken, userID).Updates(map[string]interface{}{
			"status": status,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r assetGroupInvitationRepository) DeleteAssetGroupInvitationExpired() error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableAssetGroupInvitationName).Where("expired_at < ?", jwt.GetCurrentTime()).Delete(&assets.AssetGroupInvitation{}).Error; err != nil {
			return err
		}
		return nil
	})
}
