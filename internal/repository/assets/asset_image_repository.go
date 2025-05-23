package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

// AssetImageRepository defines the interface
type AssetImageRepository interface {
	AddAssetImage(assetImage []assets.AssetImage) error
	DeleteAssetImage(assetID uint, clientID string) error
	UpdateAssetImage(assetID uint, metadata []response.AssetImageResponse, clientID string) error
	Cleanup() error
	GetAssetImageResponseByAssetID(assetID uint) (*[]response.AssetImageResponse, error)
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

func (r assetImageRepository) AddAssetImage(assetImage []assets.AssetImage) error {
	if len(assetImage) == 0 {
		return nil // No images to insert
	}

	err := r.db.Table(utils.TableAssetImageName).Create(&assetImage).Error
	if err != nil {
		log.Error().Err(err).Msg("❌ Failed to batch insert asset images")
		return err
	}

	log.Info().Int("count", len(assetImage)).Msg("✅ Successfully batch inserted asset images")
	return nil
}

// GetAssetImageResponseByAssetID retrieves asset image response by asset MaintenanceTypeID
func (r *assetImageRepository) GetAssetImageResponseByAssetID(assetID uint) (*[]response.AssetImageResponse, error) {
	var assetImages []assets.AssetImage
	err := r.db.Table(utils.TableAssetImageName).
		Where("asset_id = ?", assetID).
		Find(&assetImages).Error
	if err != nil {
		log.Error().Err(err).
			Uint("asset_id", assetID).
			Msg("❌ Failed to get asset image response by asset MaintenanceTypeID")
		return nil, err
	}

	var assetImageResponse []response.AssetImageResponse
	for _, assetImage := range assetImages {
		assetImageResponse = append(assetImageResponse, response.AssetImageResponse{
			ImageURL: assetImage.ImageURL,
		})
	}

	log.Info().
		Uint("asset_id", assetID).
		Msg("✅ Asset image response retrieved successfully")
	return &assetImageResponse, nil
}

// GetAssetImageByAssetID retrieves asset image by asset MaintenanceTypeID
func (r *assetImageRepository) GetAssetImageByAssetID(assetID uint) (*[]assets.AssetImage, error) {
	var assetImages []assets.AssetImage
	err := r.db.Table(utils.TableAssetImageName).
		Where("asset_id = ?", assetID).
		Find(&assetImages).Error
	if err != nil {
		log.Error().Err(err).
			Uint("asset_id", assetID).
			Msg("❌ Failed to get asset image by asset MaintenanceTypeID")
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
			Msg("❌ Failed to get asset images by client MaintenanceTypeID")
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

// UpdateAssetImage updates an existing asset image and logs audit
func (r *assetImageRepository) UpdateAssetImage(assetID uint, metadata []response.AssetImageResponse, clientID string) error {
	if len(metadata) == 0 {
		return nil // No images to update
	}
	// Delete existing asset images for the given assetID
	err := r.db.Unscoped().Table(utils.TableAssetImageName).
		Where("asset_id = ?", assetID).
		Delete(&assets.AssetImage{}).Error
	if err != nil {
		log.Error().Err(err).Uint("asset_id", assetID).Msg("Failed to delete existing asset images")
		return err
	}

	// Create new image records
	var newImages []assets.AssetImage
	for _, img := range metadata {
		newImages = append(newImages, assets.AssetImage{
			AssetID:      assetID,
			ImageURL:     img.ImageURL,
			UserClientID: clientID,
		})
	}

	// Insert new image metadata
	if err := r.db.Table(utils.TableAssetImageName).Create(&newImages).Error; err != nil {
		log.Error().Err(err).Uint("asset_id", assetID).Msg("Failed to create new asset images")
		return err
	}

	log.Info().
		Uint("asset_id", assetID).
		Msg("✅ Asset image updated successfully")
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
