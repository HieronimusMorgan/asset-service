package transaction

import (
	"asset-service/internal/repository/assets"
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AssetTransactionRepository interface {
	DeleteAsset(transactionID uint, clientID, fullName string) error
}

type assetTransactionRepository struct {
	db                               gorm.DB
	AssetRepository                  assets.AssetRepository
	AssetCategoryRepository          assets.AssetCategoryRepository
	AssetStatusRepository            assets.AssetStatusRepository
	AssetMaintenanceRepository       assets.AssetMaintenanceRepository
	AssetMaintenanceRecordRepository assets.AssetMaintenanceRecordRepository
	AssetAuditLogRepository          assets.AssetAuditLogRepository
}

func NewAssetTransactionRepository(db gorm.DB, AssetRepository assets.AssetRepository,
	AssetCategoryRepository assets.AssetCategoryRepository,
	AssetStatusRepository assets.AssetStatusRepository,
	AssetMaintenanceRepository assets.AssetMaintenanceRepository,
	AssetMaintenanceRecordRepository assets.AssetMaintenanceRecordRepository,
	AssetAuditLogRepository assets.AssetAuditLogRepository) AssetTransactionRepository {
	return assetTransactionRepository{
		db:                               db,
		AssetRepository:                  AssetRepository,
		AssetCategoryRepository:          AssetCategoryRepository,
		AssetStatusRepository:            AssetStatusRepository,
		AssetMaintenanceRepository:       AssetMaintenanceRepository,
		AssetMaintenanceRecordRepository: AssetMaintenanceRecordRepository,
		AssetAuditLogRepository:          AssetAuditLogRepository}
}

func (r assetTransactionRepository) DeleteAsset(transactionID uint, clientID, fullName string) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msg("🔥 Panic occurred, rolling back transaction!")
		}
	}()

	// Check if the asset exists
	checkAsset, err := r.AssetRepository.GetAssetByID(clientID, transactionID)
	if err != nil {
		tx.Rollback()
		log.Error().
			Str("method", "DeleteAsset").
			Uint("transactionID", transactionID).
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to retrieve asset")
		return err
	}
	if checkAsset.AssetID == 0 {
		tx.Rollback()
		log.Warn().
			Str("method", "DeleteAsset").
			Uint("transactionID", transactionID).
			Str("clientID", clientID).
			Msg("Asset not found, cannot proceed with deletion")
		return gorm.ErrRecordNotFound
	}

	// Check maintenance existence (skip deletion if not found)
	checkMaintenance, err := r.AssetMaintenanceRepository.GetMaintenanceByAssetID(transactionID, clientID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		log.Error().
			Str("method", "DeleteAsset").
			Uint("transactionID", transactionID).
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to check maintenance record")
		return err
	}

	// Check maintenance record existence (skip deletion if not found)
	checkMaintenanceRecord, err := r.AssetMaintenanceRecordRepository.GetMaintenanceByAssetID(transactionID, clientID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		log.Error().
			Str("method", "DeleteAsset").
			Uint("transactionID", transactionID).
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to check maintenance record")
		return err
	}

	// If maintenance record exists, delete it
	if checkMaintenanceRecord.MaintenanceRecordID != 0 {
		err = r.AssetMaintenanceRecordRepository.Delete(transactionID, fullName)
		if err != nil {
			tx.Rollback()
			log.Error().
				Str("method", "DeleteAsset").
				Uint("transactionID", transactionID).
				Str("clientID", clientID).
				Err(err).
				Msg("Failed to delete maintenance record")
			return err
		}
		log.Info().
			Uint("transactionID", transactionID).
			Str("clientID", clientID).
			Msg("✅ Maintenance record deleted")

		err = r.AssetAuditLogRepository.AfterDeleteAssetMaintenanceRecord(checkMaintenanceRecord)
	}

	// If maintenance exists, delete it
	if checkMaintenance.ID != 0 {
		err = r.AssetMaintenanceRepository.Delete(transactionID, fullName)
		if err != nil {
			tx.Rollback()
			log.Error().
				Str("method", "DeleteAsset").
				Uint("transactionID", transactionID).
				Str("clientID", clientID).
				Err(err).
				Msg("Failed to delete maintenance")
			return err
		}
		log.Info().
			Uint("transactionID", transactionID).
			Str("clientID", clientID).
			Msg("✅ Maintenance deleted")

		err = r.AssetAuditLogRepository.AfterDeleteAssetMaintenance(checkMaintenance)
	}

	// DeleteAsset the asset
	err = r.AssetRepository.DeleteAsset(transactionID, clientID)
	if err != nil {
		tx.Rollback()
		log.Error().
			Str("method", "DeleteAsset").
			Uint("transactionID", transactionID).
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to delete asset")
		return err
	}

	err = r.AssetAuditLogRepository.AfterDeleteAsset(checkAsset)

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		log.Error().
			Str("method", "DeleteAsset").
			Uint("transactionID", transactionID).
			Str("clientID", clientID).
			Err(err).
			Msg("Transaction commit failed")
		return err
	}

	log.Info().
		Uint("transactionID", transactionID).
		Str("clientID", clientID).
		Str("fullName", fullName).
		Msg("✅ Asset and related records successfully deleted")

	return nil
}

func (r assetTransactionRepository) DeleteAssetCategory(assetCategoryID uint, clientID, fullName string) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msg("🔥 Panic occurred, rolling back transaction!")
		}
	}()

	// Check if the asset category exists
	checkAssetCategory, err := r.AssetCategoryRepository.GetAssetCategoryById(assetCategoryID, clientID)
	if err != nil {
		tx.Rollback()
		log.Error().
			Str("method", "DeleteAssetCategory").
			Uint("assetCategoryID", assetCategoryID).
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to retrieve asset category")
		return err
	}
	if checkAssetCategory.AssetCategoryID == 0 {
		tx.Rollback()
		log.Warn().
			Str("method", "DeleteAssetCategory").
			Uint("assetCategoryID", assetCategoryID).
			Str("clientID", clientID).
			Msg("Asset category not found, cannot proceed with deletion")
		return gorm.ErrRecordNotFound
	}

	// Check if asset category is in use
	checkAsset, err := r.AssetRepository.GetAssetByCategoryID(assetCategoryID, clientID)
	if err != nil {
		tx.Rollback()
		log.Error().
			Str("method", "DeleteAssetCategory").
			Uint("assetCategoryID", assetCategoryID).
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to check asset category usage")
		return err
	}

	if len(checkAsset) != 0 {
		tx.Rollback()
		log.Warn().
			Str("method", "DeleteAssetCategory").
			Uint("assetCategoryID", assetCategoryID).
			Str("clientID", clientID).
			Msg("Asset category is in use, cannot proceed with deletion")
		return errors.New("asset category is in use")
	}

	// DeleteAsset the asset category
	checkAssetCategory.DeletedBy = &fullName
	err = r.AssetCategoryRepository.DeleteAssetCategory(checkAssetCategory)
	if err != nil {
		tx.Rollback()
		log.Error().
			Str("method", "DeleteAssetCategory").
			Uint("assetCategoryID", assetCategoryID).
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to delete asset category")
		return err
	}

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		log.Error().
			Str("method", "DeleteAssetCategory").
			Uint("assetCategoryID", assetCategoryID).
			Str("clientID", clientID).
			Err(err).
			Msg("Transaction commit failed")
		return err
	}

	log.Info().
		Uint("assetCategoryID", assetCategoryID).
		Str("clientID", clientID).
		Str("fullName", fullName).
		Msg("✅ Asset category successfully deleted")

	return nil
}
