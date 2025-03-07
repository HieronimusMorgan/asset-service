package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
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
func (r assetStockRepository) AddAssetStock(assetStock *assets.AssetStock) error {
	return r.db.Table(utils.TableAssetStockName).Create(&assetStock).Error
}

// DeleteAssetStock removes an asset stock and logs audit
func (r assetStockRepository) DeleteAssetStock(assetID uint, clientID string) error {
	return r.db.Table(utils.TableAssetStockName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		Updates(map[string]interface{}{"deleted_by": clientID, "deleted_at": time.Now()}).
		Delete(&assets.AssetStock{}).Error
}

// GetAssetStockResponseByAssetID retrieves asset stock by asset ID
func (r assetStockRepository) GetAssetStockResponseByAssetID(assetID uint, clientID string) (*response.AssetStockResponse, error) {
	var assetStockResponse response.AssetStockResponse
	err := r.db.Table(utils.TableAssetStockName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).Scan(&assetStockResponse).Error
	return &assetStockResponse, err
}

// GetAssetStockByAssetID retrieves asset stock by asset ID
func (r assetStockRepository) GetAssetStockByAssetID(assetID uint, clientID string) (*assets.AssetStock, error) {
	var assetStock assets.AssetStock
	err := r.db.Table(utils.TableAssetStockName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).First(&assetStock).Error
	return &assetStock, err
}

// GetAssetStock retrieves all asset stocks
func (r assetStockRepository) GetAssetStock() ([]assets.AssetStock, error) {
	var assetStocks []assets.AssetStock
	err := r.db.Table(utils.TableAssetStockName).Find(&assetStocks).Error
	return assetStocks, err
}

// GetAssetStockByClientID retrieves all asset stocks by client ID
func (r assetStockRepository) GetAssetStockByClientID(clientID string) (*[]assets.AssetStock, error) {
	var assetStocks []assets.AssetStock
	err := r.db.Table(utils.TableAssetStockName).Where("user_client_id = ?", clientID).Find(&assetStocks).Error
	return &assetStocks, err
}

func (r assetStockRepository) UpdateAssetStock(assetStock *assets.AssetStock, clientID string) error {
	tx := r.db.Begin()

	updateFields := map[string]interface{}{
		"latest_quantity": assetStock.LatestQuantity,
		"change_type":     assetStock.ChangeType,
		"updated_by":      clientID,
	}

	err := tx.Table(utils.TableAssetStockName).
		Where("asset_id = ? AND user_client_id = ?", assetStock.AssetID, clientID).
		Updates(updateFields).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
