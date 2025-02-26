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
	AddAssetImage(assetRequest []response.AssetImageResponse, assetID uint, clientID string) (interface{}, error)
	GetAssetImageByAssetID(assetID uint) (*[]assets.AssetImage, error)
	DeleteAssetImage(assetID uint, clientID string) error
	Cleanup(nats string) error
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

func (s assetImageService) AddAssetImage(assetRequest []response.AssetImageResponse, assetID uint, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return nil, err
	}

	if len(assetRequest) != 0 {
		for _, image := range assetRequest {
			var assetImage = &assets.AssetImage{
				AssetID:   assetID,
				ImageURL:  image.ImageURL,
				FileSize:  image.FileSize,
				FileType:  image.FileType,
				CreatedBy: data.ClientID,
				UpdatedBy: data.ClientID,
			}
			err = s.AssetImageRepository.AddAssetImage(assetImage)
		}

		if err != nil {
			log.Error().
				Str("key", "AddAssetImage").
				Str("clientID", clientID).
				Err(err).
				Msg("Failed to add asset image")
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
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

// ImageDeleteRequest represents the deletion request sent to NATS
type ImageDeleteRequest struct {
	ClientID string   `json:"client_id"`
	Images   []string `json:"images"`
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
		}
		if err = requestImageDeletion(nats, asset.UserClientID, images); err != nil {
			log.Error().Str("key", "RequestImageDeletion").Err(err).Msg("Failed to request image deletion")
		}
	}

	return nil
}

// Sends a delete request to `cdn-service` via NATS
func requestImageDeletion(natsAPI string, clientID string, images []string) error {
	nc, err := nats.Connect(natsAPI)
	if err != nil {
		return err
	}
	defer nc.Close()

	data, _ := json.Marshal(ImageDeleteRequest{ClientID: clientID, Images: images})

	return nc.Publish("asset.image.delete", data)
}
