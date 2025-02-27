package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

// AssetImageRepository defines the interface
type AssetImageRepository interface {
	AddAssetImage(assetImage *assets.AssetImage) error
	DeleteAssetImage(assetID uint, clientID string) error
	Cleanup() error
	GetAssetImageByAssetID(assetID uint) (*[]assets.AssetImage, error)
	GetAssetImage() ([]assets.AssetImage, error)
	GetAssetImageByClientID(clientID string) (*[]assets.AssetImage, error)
}

// assetImageRepository implementation
type assetImageRepository struct {
	db gorm.DB
}

// NewAssetImageRepository initializes the repository
func NewAssetImageRepository(db gorm.DB) AssetImageRepository {
	return &assetImageRepository{db: db}
}

// AddAssetImage inserts a new asset image and logs audit

func (r assetImageRepository) AddAssetImage(assetImage *assets.AssetImage) error {
	if assetImage == nil {
		return errors.New("assets cannot be nil")
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table(utils.TableAssetImageName).Create(&assetImage).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create assets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAssetImageByAssetID retrieves asset image by asset ID
func (r *assetImageRepository) GetAssetImageByAssetID(assetID uint) (*[]assets.AssetImage, error) {
	var assetImages []assets.AssetImage
	err := r.db.Table(utils.TableAssetImageName).
		Where("asset_id = ?", assetID).
		Find(&assetImages).Error
	if err != nil {
		log.Error().Err(err).
			Uint("asset_id", assetID).
			Msg("❌ Failed to get asset image by asset ID")
		return nil, err
	}

	log.Info().
		Uint("asset_id", assetID).
		Msg("✅ Asset image retrieved successfully")
	return &assetImages, nil
}

// GetAssetImage retrieves all asset images
func (r *assetImageRepository) GetAssetImage() ([]assets.AssetImage, error) {
	var assetImages []assets.AssetImage
	err := r.db.Table(utils.TableAssetImageName).
		Find(&assetImages).Error
	if err != nil {
		log.Error().Err(err).Msg("❌ Failed to get asset images")
		return nil, err
	}

	log.Info().Msg("✅ Asset images retrieved successfully")
	return assetImages, nil
}

func (r *assetImageRepository) GetAssetImageByClientID(clientID string) (*[]assets.AssetImage, error) {
	var assetImages []assets.AssetImage
	err := r.db.Table(utils.TableAssetImageName).
		Where("user_client_id = ?", clientID).
		Find(&assetImages).Error
	if err != nil {
		log.Error().Err(err).
			Str("client_id", clientID).
			Msg("❌ Failed to get asset images by client ID")
		return nil, err
	}

	log.Info().
		Str("client_id", clientID).
		Msg("✅ Asset images retrieved successfully")
	return &assetImages, nil
}

// DeleteAssetImage removes an existing asset image and logs audit
func (r *assetImageRepository) DeleteAssetImage(assetID uint, clientID string) error {
	err := r.db.Table(utils.TableAssetImageName).
		Where("asset_id = ?", assetID).
		Updates(map[string]interface{}{"deleted_by": clientID, "deleted_at": time.Now()}).
		Delete(&assets.Asset{}).Error
	if err != nil {
		log.Error().Err(err).
			Uint("asset_id", assetID).
			Msg("❌ Failed to delete asset image")
		return err
	}

	log.Info().
		Uint("asset_id", assetID).
		Msg("✅ Asset image deleted successfully")
	return nil
}

// Cleanup removes all asset images
func (r *assetImageRepository) Cleanup() error {
	err := r.db.Table(utils.TableAssetImageName).
		Delete(&assets.AssetImage{}).Error
	if err != nil {
		log.Error().Err(err).Msg("❌ Failed to cleanup asset images")
		return err
	}

	log.Info().Msg("✅ Asset images cleaned up successfully")
	return nil
}
