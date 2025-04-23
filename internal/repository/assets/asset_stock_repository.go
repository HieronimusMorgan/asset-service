package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

// AssetStockRepository defines the interface
type AssetStockRepository interface {
	AddAssetStock(assetStock *assets.AssetStock) error
	DeleteAssetStock(assetID uint, clientID string) error
	GetAssetStockResponseByAssetID(assetID uint, clientID string) (*response.AssetStockResponse, error)
	GetAssetStockByAssetID(assetID uint, clientID string) (*assets.AssetStock, error)
	GetAssetStock() ([]assets.AssetStock, error)
	GetAssetStockByClientID(clientID string) (*[]assets.AssetStock, error)
	UpdateAssetStock(assetStock *assets.AssetStock, clientID string) error
	GetAssetStockByAssetIDAndAssetGroupID(assetID, assetGroupID uint) (*assets.AssetStock, error)
	UpdateAssetStockByAssetGroupID(assetStock *assets.AssetStock, assetGroupID uint, clientID string) error
}

// assetStockRepository implementation
type assetStockRepository struct {
	db gorm.DB
}

// NewAssetStockRepository initializes the repository
func NewAssetStockRepository(db gorm.DB) AssetStockRepository {
	return &assetStockRepository{db: db}
}

// AddAssetStock inserts a new asset stock and logs audit
func (r *assetStockRepository) AddAssetStock(assetStock *assets.AssetStock) error {
	return r.db.Table(utils.TableAssetStockName).Create(&assetStock).Error
}

// DeleteAssetStock removes an asset stock and logs audit
func (r *assetStockRepository) DeleteAssetStock(assetID uint, clientID string) error {
	return r.db.Table(utils.TableAssetStockName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		Updates(map[string]interface{}{"deleted_by": clientID, "deleted_at": time.Now()}).
		Delete(&assets.AssetStock{}).Error
}

// GetAssetStockResponseByAssetID retrieves asset stock by asset MaintenanceTypeID
func (r *assetStockRepository) GetAssetStockResponseByAssetID(assetID uint, clientID string) (*response.AssetStockResponse, error) {
	var assetStockResponse response.AssetStockResponse
	err := r.db.Table(utils.TableAssetStockName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).Scan(&assetStockResponse).Error
	return &assetStockResponse, err
}

// GetAssetStockByAssetID retrieves asset stock by asset MaintenanceTypeID
func (r *assetStockRepository) GetAssetStockByAssetID(assetID uint, clientID string) (*assets.AssetStock, error) {
	var assetStock assets.AssetStock
	err := r.db.Table(utils.TableAssetStockName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).First(&assetStock).Error
	return &assetStock, err
}

// GetAssetStock retrieves all asset stocks
func (r *assetStockRepository) GetAssetStock() ([]assets.AssetStock, error) {
	var assetStocks []assets.AssetStock
	err := r.db.Table(utils.TableAssetStockName).Find(&assetStocks).Error
	return assetStocks, err
}

// GetAssetStockByClientID retrieves all asset stocks by client MaintenanceTypeID
func (r *assetStockRepository) GetAssetStockByClientID(clientID string) (*[]assets.AssetStock, error) {
	var assetStocks []assets.AssetStock
	err := r.db.Table(utils.TableAssetStockName).Where("user_client_id = ?", clientID).Find(&assetStocks).Error
	return &assetStocks, err
}

func (r *assetStockRepository) UpdateAssetStock(assetStock *assets.AssetStock, clientID string) error {
	tx := r.db.Begin()

	var existingStock assets.AssetStock
	if err := tx.Table(utils.TableAssetStockName).
		Where("asset_id = ? AND user_client_id = ?", assetStock.AssetID, clientID).
		First(&existingStock).Error; err != nil {
		tx.Rollback()
		return err
	}

	previousQuantity := existingStock.LatestQuantity
	newQuantity := previousQuantity

	switch assetStock.ChangeType {
	case "INCREASE":
		newQuantity += assetStock.Quantity
	case "DECREASE":
		if previousQuantity < assetStock.Quantity {
			tx.Rollback()
			return errors.New("not enough stock available")
		}
		newQuantity -= assetStock.Quantity
	case "ADJUSTMENT":
		newQuantity = assetStock.Quantity
	default:
		tx.Rollback()
		return errors.New("invalid stock change type")
	}

	// If quantity has not changed, do not insert into history
	if newQuantity == previousQuantity {
		tx.Rollback()
		return errors.New("no stock change detected")
	}

	// Update asset_stock table
	updateFields := map[string]interface{}{
		"latest_quantity": newQuantity,
		"quantity":        newQuantity,
		"change_type":     assetStock.ChangeType,
		"reason":          assetStock.Reason,
		"updated_by":      clientID,
		"updated_at":      time.Now(),
	}

	if err := tx.Table(utils.TableAssetStockName).
		Where("asset_id = ? AND user_client_id = ?", assetStock.AssetID, clientID).
		Updates(updateFields).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insert into asset_stock_history
	stockHistory := assets.AssetStockHistory{
		AssetID:          assetStock.AssetID,
		UserClientID:     clientID,
		StockID:          existingStock.StockID,
		ChangeType:       assetStock.ChangeType,
		PreviousQuantity: previousQuantity,
		NewQuantity:      newQuantity,
		QuantityChanged:  abs(newQuantity - previousQuantity), // Ensure this is always > 0
		Reason:           assetStock.Reason,
		CreatedBy:        clientID,
		CreatedAt:        time.Now(),
	}

	if err := tx.Table(utils.TableAssetStockHistoryName).Create(&stockHistory).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *assetStockRepository) GetAssetStockByAssetIDAndAssetGroupID(assetID, assetGroupID uint) (*assets.AssetStock, error) {
	var assetStock assets.AssetStock
	query := `
		SELECT stock.*
		FROM "asset_stock" stock
		LEFT JOIN asset a ON stock.asset_id = a.asset_id
		LEFT JOIN "asset_group_asset" aga ON a.asset_id = aga.asset_id
		WHERE a.asset_id = ? AND aga.asset_group_id = ?;
	`
	err := r.db.Raw(query, assetID, assetGroupID).First(&assetStock).Error
	if err != nil {
		return nil, err
	}
	return &assetStock, nil
}

func (r *assetStockRepository) UpdateAssetStockByAssetGroupID(assetStock *assets.AssetStock, assetGroupID uint, clientID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingStock assets.AssetStock
		query := `
		SELECT stock.*
		FROM "asset_stock" stock
		LEFT JOIN asset a ON stock.asset_id = a.asset_id
		LEFT JOIN "asset_group_asset" aga ON a.asset_id = aga.asset_id
		WHERE a.asset_id = ? AND aga.asset_group_id = ?;
	`
		err := r.db.Raw(query, assetStock.AssetID, assetGroupID).First(&existingStock).Error
		if err != nil {
			return err
		}

		log.Info().Interface("existingStock", existingStock).Msg("Existing stock retrieved successfully")

		previousQuantity := existingStock.LatestQuantity
		newQuantity := previousQuantity

		switch assetStock.ChangeType {
		case "INCREASE":
			newQuantity += assetStock.Quantity
		case "DECREASE":
			if previousQuantity < assetStock.Quantity {
				tx.Rollback()
				return errors.New("not enough stock available")
			}
			newQuantity -= assetStock.Quantity
		case "ADJUSTMENT":
			newQuantity = assetStock.Quantity
		default:
			tx.Rollback()
			return errors.New("invalid stock change type")
		}

		// If quantity has not changed, do not insert into history
		if newQuantity == previousQuantity {
			tx.Rollback()
			return errors.New("no stock change detected")
		}

		if err := tx.Table(utils.TableAssetStockName).
			Where("asset_id = ? AND stock_id = ?", assetStock.AssetID, existingStock.StockID).
			Updates(map[string]interface{}{
				"latest_quantity": newQuantity,
				"quantity":        newQuantity,
				"change_type":     assetStock.ChangeType,
				"reason":          utils.NilIfEmpty(*assetStock.Reason),
				"updated_by":      clientID,
				"updated_at":      time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Insert into asset_stock_history
		stockHistory := assets.AssetStockHistory{
			AssetID:          assetStock.AssetID,
			UserClientID:     clientID,
			StockID:          existingStock.StockID,
			ChangeType:       assetStock.ChangeType,
			PreviousQuantity: previousQuantity,
			NewQuantity:      newQuantity,
			QuantityChanged:  abs(newQuantity - previousQuantity), // Ensure this is always > 0
			Reason:           assetStock.Reason,
			CreatedBy:        clientID,
			CreatedAt:        time.Now(),
		}

		if err = tx.Table(utils.TableAssetStockHistoryName).Create(&stockHistory).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

// Helper function to ensure absolute values
func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
