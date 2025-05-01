package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repo "asset-service/internal/repository/assets"
	repouser "asset-service/internal/repository/users"
	"asset-service/internal/utils"
	"asset-service/internal/utils/redis"
	"errors"
	"github.com/rs/zerolog/log"
)

type AssetWishlistService interface {
	AddAssetWishlist(assetRequest *request.AssetWishlistRequest, clientID string) (interface{}, error)
	GetListAssetWishlist(clientID string, size, index int) (interface{}, int64, error)
	GetAssetWishlistByID(id uint, clientID string) (interface{}, error)
	UpdateAssetWishlist(id uint, req *request.AssetWishlistRequest, clientID string) (interface{}, error)
	DeleteAssetWishlist(id uint, clientID string) error
	AddAssetWishlistToAsset(id uint, req *request.AssetRequest, metadata []response.AssetImageResponse, clientID string, headerID string) (interface{}, error)
}

type assetWishlistService struct {
	UserRepository             repouser.UserRepository
	AssetWishlistRepository    repo.AssetWishlistRepository
	AssetRepository            repo.AssetRepository
	AssetCategoryRepository    repo.AssetCategoryRepository
	AssetStatusRepository      repo.AssetStatusRepository
	AssetImageRepository       repo.AssetImageRepository
	AuditLogRepository         repo.AssetAuditLogRepository
	AssetGroupMemberRepository repo.AssetGroupMemberRepository
	AssetGroupAssetRepository  repo.AssetGroupAssetRepository
	Redis                      redis.RedisService
}

func NewAssetWishlistService(
	UserRepository repouser.UserRepository,
	assetWishlistRepository repo.AssetWishlistRepository,
	AssetCategoryRepository repo.AssetCategoryRepository,
	AssetStatusRepository repo.AssetStatusRepository,
	AssetImageRepository repo.AssetImageRepository,
	assetRepository repo.AssetRepository,
	AuditLogRepository repo.AssetAuditLogRepository,
	AssetGroupMemberRepository repo.AssetGroupMemberRepository,
	AssetGroupAssetRepository repo.AssetGroupAssetRepository,
	redis redis.RedisService) AssetWishlistService {
	return assetWishlistService{
		UserRepository:             UserRepository,
		AssetWishlistRepository:    assetWishlistRepository,
		AssetRepository:            assetRepository,
		AssetCategoryRepository:    AssetCategoryRepository,
		AssetStatusRepository:      AssetStatusRepository,
		AssetImageRepository:       AssetImageRepository,
		AuditLogRepository:         AuditLogRepository,
		AssetGroupMemberRepository: AssetGroupMemberRepository,
		AssetGroupAssetRepository:  AssetGroupAssetRepository,
		Redis:                      redis,
	}
}

func (s assetWishlistService) AddAssetWishlist(assetRequest *request.AssetWishlistRequest, clientID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	check, err := s.AssetWishlistRepository.AssetWishlistNameExists(assetRequest.AssetName, data.ClientID)
	if err != nil {
		log.Error().
			Str("key", "AssetNameExists").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to check assetWishlist name exists")
		return nil, err
	}
	if check {
		log.Error().
			Str("key", "AssetNameExists").
			Str("clientID", clientID).
			Err(err).
			Msg("Asset name already exists")
		return nil, errors.New("assets already exists")
	}

	var assetWishlist = assets.AssetWishlist{
		UserClientID:  data.ClientID,
		AssetName:     assetRequest.AssetName,
		SerialNumber:  assetRequest.SerialNumber,
		Barcode:       assetRequest.Barcode,
		CategoryID:    assetRequest.CategoryID,
		StatusID:      assetRequest.StatusID,
		PriorityLevel: assetRequest.PriorityLevel,
		PriceEstimate: assetRequest.PriceEstimate,
		Notes:         assetRequest.Notes,
		CreatedBy:     &data.ClientID,
		UpdatedBy:     &data.ClientID,
	}

	err = s.AssetWishlistRepository.AddAssetWishlist(&assetWishlist)
	if err != nil {
		log.Error().
			Str("key", "AddAssetWishlist").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to add assetWishlist wishlist")
		return nil, err
	}

	assetResult, err := s.AssetWishlistRepository.GetAssetWishlistResponseByID(clientID, assetWishlist.WishlistID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get assetWishlist by MaintenanceTypeID")
		return nil, err
	}

	return assetResult, nil
}

func (s assetWishlistService) GetListAssetWishlist(clientID string, size, index int) (interface{}, int64, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, 0, err
	}

	total, err := s.AssetWishlistRepository.GetListAssetWishlistCount(data.ClientID)
	if err != nil {
		log.Error().
			Str("key", "GetListAssetWishlistCount").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset wishlist count")
		return nil, 0, err
	}

	assetWishlist, err := s.AssetWishlistRepository.GetListAssetWishlist(data.ClientID, size, index)
	if err != nil {
		log.Error().
			Str("key", "GetListAssetWishlist").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset wishlist")
		return nil, total, err
	}

	return assetWishlist, total, nil
}

func (s assetWishlistService) GetAssetWishlistByID(id uint, clientID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	assetWishlist, err := s.AssetWishlistRepository.GetAssetWishlistResponseByID(data.ClientID, id)
	if err != nil {
		log.Error().
			Str("key", "GetAssetWishlistResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset wishlist by ID")
		return nil, err
	}

	return assetWishlist, nil
}

func (s assetWishlistService) UpdateAssetWishlist(id uint, req *request.AssetWishlistRequest, clientID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	check, err := s.AssetWishlistRepository.AssetWishlistNameExists(req.AssetName, data.ClientID)
	if err != nil {
		log.Error().
			Str("key", "AssetNameExists").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to check assetWishlist name exists")
		return nil, err
	}
	if check {
		log.Error().
			Str("key", "AssetNameExists").
			Str("clientID", clientID).
			Err(err).
			Msg("Asset name already exists")
		return nil, errors.New("assets already exists")
	}

	var assetWishlist = assets.AssetWishlist{
		WishlistID:    id,
		UserClientID:  data.ClientID,
		AssetName:     req.AssetName,
		SerialNumber:  req.SerialNumber,
		Barcode:       req.Barcode,
		CategoryID:    req.CategoryID,
		StatusID:      req.StatusID,
		PriorityLevel: req.PriorityLevel,
		PriceEstimate: req.PriceEstimate,
		Notes:         req.Notes,
		UpdatedBy:     &data.ClientID,
	}

	err = s.AssetWishlistRepository.UpdateAssetWishlist(&assetWishlist)
	if err != nil {
		log.Error().
			Str("key", "UpdateAssetWishlist").
			Uint("id", id).
			Err(err).
			Msg("Failed to update asset wishlist")
		return nil, err
	}

	return assetWishlist, nil
}

func (s assetWishlistService) DeleteAssetWishlist(id uint, clientID string) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return err
	}

	err = s.AssetWishlistRepository.DeleteAssetWishlist(id, data.ClientID)
	if err != nil {
		log.Error().
			Str("key", "DeleteAssetWishlist").
			Uint("id", id).
			Err(err).
			Msg("Failed to delete asset wishlist")
		return err
	}

	return nil
}

func (s assetWishlistService) AddAssetWishlistToAsset(id uint, assetRequest *request.AssetRequest, images []response.AssetImageResponse, clientID string, headerID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logError("GetUserByClientID", clientID, err, "Failed to get user by client MaintenanceTypeID")
	}

	assetWishlist, err := s.AssetWishlistRepository.GetAssetWishlistByID(data.ClientID, id)
	if err != nil {
		return logError("GetAssetWishlistByID", clientID, err, "Failed to get asset wishlist by MaintenanceTypeID")
	}
	if assetWishlist == nil {
		log.Error().
			Str("key", "GetAssetWishlistByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Asset wishlist not found")
		return nil, errors.New("asset wishlist not found")
	}

	if exists, err := s.AssetRepository.AssetNameExists(assetRequest.Name, clientID); err != nil {
		return logError("AssetNameExists", clientID, err, "Failed to check asset name exists")
	} else if exists {
		return logError("AssetNameExists", clientID, errors.New("assets already exists"), "Asset name already exists")
	}

	if _, err := s.AssetCategoryRepository.GetAssetCategoryById(uint(assetRequest.CategoryID), clientID); err != nil {
		return logError("GetAssetCategoryById", clientID, errors.New("category not found"), "Failed to get asset category by MaintenanceTypeID")
	}

	if _, err := s.AssetStatusRepository.GetAssetStatusByID(uint(assetRequest.StatusID)); err != nil {
		return logError("GetAssetStatusByID", clientID, errors.New("status not found"), "Failed to get asset status by MaintenanceTypeID")
	}

	purchaseDate, _ := utils.ParseOptionalDate(assetRequest.PurchaseDate)
	expiryDate, _ := utils.ParseOptionalDate(assetRequest.ExpiryDate)
	warrantyExpiry, _ := utils.ParseOptionalDate(assetRequest.WarrantyExpiry)

	assetRequest.ConvertAssetRequestEmptyToNil()

	asset := &assets.Asset{
		UserClientID:       clientID,
		SerialNumber:       assetRequest.SerialNumber,
		Name:               assetRequest.Name,
		Description:        assetRequest.Description,
		Barcode:            assetRequest.Barcode,
		CategoryID:         assetRequest.CategoryID,
		StatusID:           assetRequest.StatusID,
		PurchaseDate:       purchaseDate,
		ExpiryDate:         expiryDate,
		WarrantyExpiryDate: warrantyExpiry,
		Price:              assetRequest.Price,
		Stock:              assetRequest.Stock,
		Notes:              assetRequest.Notes,
		CreatedBy:          &data.ClientID,
		UpdatedBy:          &data.ClientID,
	}

	if err := s.AssetRepository.AddAssetFromWishlist(asset, assetWishlist, images); err != nil {
		return logError("AddAsset", clientID, err, "Failed to add asset")
	}

	// check user have asset group
	assetGroupMember, _ := s.AssetGroupMemberRepository.GetAssetGroupMemberByUserID(user.UserID)
	if assetGroupMember != nil && assetGroupMember.UserID != 0 {
		assetGroupAsset := &assets.AssetGroupAsset{
			AssetID:      asset.AssetID,
			AssetGroupID: assetGroupMember.AssetGroupID,
			UserID:       user.UserID,
			CreatedBy:    &data.ClientID,
		}

		if err := s.AssetGroupAssetRepository.AddAssetGroupAsset(assetGroupAsset); err != nil {
			return logError("AddAssetGroupAsset", clientID, err, "Failed to add asset group asset")
		}
	}

	assetImage, err := s.AssetImageRepository.GetAssetImageResponseByAssetID(asset.AssetID)
	if err != nil {
		return logError("GetAssetImageResponseByAssetID", clientID, err, "Failed to get asset image by MaintenanceTypeID")
	}

	result, err := s.AssetRepository.GetAssetResponseByID(clientID, asset.AssetID)
	if err != nil {
		return logError("GetAssetResponseByID", clientID, err, "Failed to get asset by MaintenanceTypeID")
	}

	result.Images = *assetImage

	log.Info().Str("key", "GetAssetResponseByID").Str("clientID", clientID).Fields(asset).Msg("Success to get asset by MaintenanceTypeID")
	return result, nil
}
