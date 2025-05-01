package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	nt "asset-service/internal/utils/nats"
	"asset-service/internal/utils/redis"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

type AssetImageService interface {
	AddAssetImage(assetRequest []response.AssetImageResponse, assetID uint, clientID string) error
	GetAssetImageByAssetID(assetID uint) (*[]assets.AssetImage, error)
	DeleteAssetImage(assetID uint, clientID string) error
	Cleanup() error
	CleanupUnusedImages() error
}

type assetImageService struct {
	AssetImageRepository repository.AssetImageRepository
	AssetRepository      repository.AssetRepository
	Redis                redis.RedisService
	NatsService          nt.Service
}

func NewAssetImageService(assetImageRepository repository.AssetImageRepository, assetRepository repository.AssetRepository, redis redis.RedisService, natsService nt.Service) AssetImageService {
	return &assetImageService{
		AssetImageRepository: assetImageRepository,
		AssetRepository:      assetRepository,
		Redis:                redis,
		NatsService:          natsService,
	}
}

func (s assetImageService) AddAssetImage(assetRequest []response.AssetImageResponse, assetID uint, clientID string) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return err
	}

	if len(assetRequest) != 0 {
		var assetImages []assets.AssetImage
		for _, image := range assetRequest {
			assetImages = append(assetImages, assets.AssetImage{
				UserClientID: clientID,
				AssetID:      assetID,
				ImageURL:     image.ImageURL,
				CreatedBy:    &data.ClientID,
				UpdatedBy:    &data.ClientID,
			})
		}
		if err := s.AssetImageRepository.AddAssetImage(assetImages); err != nil {
			log.Error().
				Str("key", "AddAssetImage").
				Str("clientID", clientID).
				Err(err).
				Msg("Failed to add asset image")
			return err
		}
	}
	return nil
}

func (s assetImageService) GetAssetImageByAssetID(assetID uint) (*[]assets.AssetImage, error) {
	assetImage, err := s.AssetImageRepository.GetAssetImageByAssetID(assetID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetImageByAssetID").
			Uint("asset_id", assetID).
			Err(err).
			Msg("Failed to get asset image")
		return nil, err
	}

	return assetImage, nil
}

func (s assetImageService) DeleteAssetImage(assetID uint, clientID string) error {
	err := s.AssetImageRepository.DeleteAssetImage(assetID, clientID)
	if err != nil {
		log.Error().
			Str("key", "DeleteAssetImage").
			Uint("asset_id", assetID).
			Err(err).
			Msg("Failed to delete asset image")
		return err
	}

	return nil
}

// Cleanup removes images of deleted assets
func (s assetImageService) Cleanup() error {
	deleted, err := s.AssetRepository.GetAssetDeleted()
	if err != nil {
		log.Error().Str("key", "GetAssetDeleted").Err(err).Msg("Failed to get deleted assets")
		return err
	}

	if len(deleted) == 0 {
		log.Info().Msg("✅ No deleted assets found, cleanup not needed.")
		return nil
	}

	for _, asset := range deleted {
		// Fetch images associated with the asset
		assetImages, err := s.AssetImageRepository.GetAssetImageByAssetID(asset.AssetID)
		if err != nil {
			log.Error().Str("key", "GetAssetImageByAssetID").Uint("asset_id", asset.AssetID).Err(err).Msg("Failed to get asset images")
			return err
		}

		// Collect image filenames
		var images []string
		for _, img := range *assetImages {
			images = append(images, filepath.Base(img.ImageURL))
		}

		// **VALIDATION: If no images, return an error**
		if len(images) == 0 {
			errMsg := fmt.Sprintf("❌ No images found for asset MaintenanceTypeID %d, skipping deletion", asset.AssetID)
			log.Warn().Str("client_id", asset.UserClientID).Msg(errMsg)
			return errors.New(errMsg)
		}

		// Delete asset image metadata from database
		if err = s.DeleteAssetImage(asset.AssetID, asset.UserClientID); err != nil {
			log.Error().Str("key", "DeleteAssetImage").Uint("asset_id", asset.AssetID).Err(err).Msg("Failed to delete asset images")
			return err
		}

		// Log the number of images being deleted
		log.Info().Str("client_id", asset.UserClientID).Msgf("🗑️ Images to be deleted: %d", len(images))

		// Request physical file deletion via NATS
		if err = s.NatsService.RequestImageDeletion(asset.UserClientID, images); err != nil {
			log.Error().Str("key", "RequestImageDeletion").Err(err).Msg("❌ Failed to request image deletion via NATS")
			return err
		}
	}

	log.Info().Msg("✅ Cleanup process completed successfully.")
	return nil
}

func (s assetImageService) CleanupUnusedImages() error {
	images, err := s.AssetImageRepository.GetAssetImage()
	if err != nil {
		log.Error().Str("key", "GetUnusedImages").Err(err).Msg("Failed to get unused images")
		return err
	}
	var assetImages []assets.ImageDeleteRequest
	var clientID string
	var url []string
	for _, img := range images {
		clientID = img.UserClientID
		imagesAsset, err := s.AssetImageRepository.GetAssetImageByClientID(clientID)
		if err != nil {
			log.Error().Str("key", "GetAssetImageByClientID").Str("client_id", clientID).Err(err).Msg("Failed to get asset images")
		}
		if imagesAsset == nil {
			continue
		}
		for _, image := range *imagesAsset {
			url = append(url, filepath.Base(image.ImageURL))
		}
		assetImages = append(assetImages, assets.ImageDeleteRequest{
			ClientID: clientID,
			Images:   url,
		})
	}
	if len(assetImages) == 0 {
		log.Info().Msg("No unused images found")
		return nil
	} else {
		if err = s.NatsService.RequestImageUsage(assetImages); err != nil {
			log.Error().Str("key", "RequestImageDeletion").Err(err).Msg("Failed to request image deletion")
		}
	}
	return nil
}
