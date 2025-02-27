package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

type AssetImageService interface {
	AddAssetImage(assetRequest []response.AssetImageResponse, assetID uint, clientID string) error
	GetAssetImageByAssetID(assetID uint) (*[]assets.AssetImage, error)
	DeleteAssetImage(assetID uint, clientID string) error
	Cleanup(nats string) error
	CleanupUnusedImages(nats string) error
}

type assetImageService struct {
	AssetImageRepository repository.AssetImageRepository
	AssetRepository      repository.AssetRepository
	Redis                utils.RedisService
}

func NewAssetImageService(
	assetImageRepository repository.AssetImageRepository,
	assetRepository repository.AssetRepository,
	redis utils.RedisService) AssetImageService {
	return &assetImageService{
		AssetImageRepository: assetImageRepository,
		AssetRepository:      assetRepository,
		Redis:                redis,
	}
}

func (s assetImageService) AddAssetImage(assetRequest []response.AssetImageResponse, assetID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return err
	}

	if len(assetRequest) != 0 {
		for _, image := range assetRequest {
			var assetImage = &assets.AssetImage{
				UserClientID: clientID,
				AssetID:      assetID,
				ImageURL:     image.ImageURL,
				FileSize:     image.FileSize,
				FileType:     image.FileType,
				CreatedBy:    data.ClientID,
				UpdatedBy:    data.ClientID,
			}
			err = s.AssetImageRepository.AddAssetImage(assetImage)
		}

		if err != nil {
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
func (s assetImageService) Cleanup(nats string) error {
	deleted, err := s.AssetRepository.GetAssetDeleted()
	if err != nil {
		log.Error().Str("key", "GetAssetDeleted").Err(err).Msg("Failed to get deleted assets")
		return err
	}

	for _, asset := range deleted {
		assetImages, err := s.AssetImageRepository.GetAssetImageByAssetID(asset.AssetID)
		if err != nil {
			log.Error().Str("key", "GetAssetImageByAssetID").Uint("asset_id", asset.AssetID).Err(err).Msg("Failed to get asset images")
			return err
		}

		var images []string
		for _, img := range *assetImages {
			images = append(images, filepath.Base(img.ImageURL))
		}

		if err = s.DeleteAssetImage(asset.AssetID, asset.UserClientID); err != nil {
			log.Error().Str("key", "DeleteAssetImage").Uint("asset_id", asset.AssetID).Err(err).Msg("Failed to delete asset images")
			return err
		}

		if len(images) == 0 {
			continue
		} else {
			if err = requestImageDeletion(nats, asset.UserClientID, images); err != nil {
				log.Error().Str("key", "RequestImageDeletion").Err(err).Msg("Failed to request image deletion")
			}
		}
	}

	return nil
}

func (s assetImageService) CleanupUnusedImages(nats string) error {
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
		if err = requestImageUsage(nats, assetImages); err != nil {
			log.Error().Str("key", "RequestImageDeletion").Err(err).Msg("Failed to request image deletion")
		}
	}
	return nil
}

func requestImageDeletion(natsAPI string, clientID string, images []string) error {
	nc, err := nats.Connect(natsAPI)
	if err != nil {
		return err
	}
	defer nc.Close()

	data, _ := json.Marshal(assets.ImageDeleteRequest{ClientID: clientID, Images: images})

	return nc.Publish("asset.image.delete", data)
}

func requestImageUsage(natsAPI string, images []assets.ImageDeleteRequest) error {
	nc, err := nats.Connect(natsAPI)
	if err != nil {
		return err
	}
	defer nc.Close()

	data, _ := json.Marshal(images)

	return nc.Publish("asset.image.usage", data)
}
