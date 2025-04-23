package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type AssetRepository interface {
	AddAsset(asset *assets.Asset, images []response.AssetImageResponse) error
	GetAssetByNameAndClientID(name string, clientID string) (*assets.Asset, error)
	AssetNameExists(name string, clientID string) (bool, error)
	GetAsset(assetID uint, clientID string) (*assets.Asset, error)
	GetAssetByAssetGroupID(assetID, assetGroupID uint) (*assets.Asset, error)
	GetCountAsset(clientID string) (int64, error)
	GetListAssets(clientID string, index int, size int) ([]response.AssetResponse, error)
	GetListAssetsByAssetGroup(clientID string, assetGroupID uint) ([]response.AssetResponse, error)
	GetAssetResponseByID(clientID string, id uint) (*response.AssetResponse, error)
	GetAssetByID(clientID string, id uint) (*assets.Asset, error)
	UpdateAsset(asset *assets.Asset, clientID string) error
	UpdateMaintenanceDateAsset(assetID uint, maintenanceDate *time.Time, clientID string) error
	UpdateAssetStatus(assetID uint, statusID uint, clientID string) (*assets.Asset, error)
	UpdateAssetCategory(assetID uint, categoryID uint, clientID string) (*assets.Asset, error)
	DeleteAsset(id uint, clientID string) error
	GetAssetByIDForMaintenance(id uint, clientID string) (*assets.Asset, error)
	GetAssetByCategoryID(assetCategoryID uint, clientID string) ([]assets.Asset, error)
	GetAssetDeleted() ([]assets.Asset, error)
}

type assetRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
}

func NewAssetRepository(db gorm.DB, audit AssetAuditLogRepository) AssetRepository {
	return assetRepository{db: db, audit: audit}
}

func (r assetRepository) AddAsset(asset *assets.Asset, images []response.AssetImageResponse) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableAssetName).Create(&asset).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create assets: %w", err)
		}

		if len(images) > 0 {
			var assetImages []assets.AssetImage
			for _, image := range images {
				assetImages := append(assetImages, assets.AssetImage{
					UserClientID: asset.UserClientID,
					AssetID:      asset.AssetID,
					ImageURL:     image.ImageURL,
					CreatedBy:    asset.UserClientID,
					UpdatedBy:    asset.UserClientID,
				})
				if err := tx.Table(utils.TableAssetImageName).Create(&assetImages).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to create asset images: %w", err)
				}
			}
		}

		assetStock := &assets.AssetStock{
			AssetID:         asset.AssetID,
			UserClientID:    asset.UserClientID,
			InitialQuantity: asset.Stock,
			LatestQuantity:  asset.Stock,
			ChangeType:      "INCREASE",
			Quantity:        asset.Stock,
			Reason:          nil,
			CreatedBy:       asset.UserClientID,
		}
		if err := tx.Table(utils.TableAssetStockName).Create(&assetStock).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create asset stock: %w", err)
		}

		return nil
	})
}

func (r assetRepository) GetAssetByNameAndClientID(name string, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("name = ? AND user_client_id = ?", name, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetRepository) AssetNameExists(name string, clientID string) (bool, error) {
	var count int64
	err := r.db.Table(utils.TableAssetName).
		Where("name = ? AND user_client_id = ?", name, clientID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r assetRepository) GetAsset(assetID uint, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.db.Table(utils.TableAssetName).
		Where("asset_id = ? AND user_client_id = ?", assetID, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetRepository) GetAssetByAssetGroupID(assetID, assetGroupID uint) (*assets.Asset, error) {
	var asset assets.Asset
	query := `
		SELECT a.*
		FROM "asset" a
		LEFT JOIN "asset_group_asset" aga ON a.asset_id = aga.asset_id
		WHERE a.asset_id = ? AND aga.asset_group_id = ?;
	`
	err := r.db.Raw(query, assetID, assetGroupID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetRepository) GetCountAsset(clientID string) (int64, error) {
	var count int64
	err := r.db.Table(utils.TableAssetName).
		Where("user_client_id = ? AND deleted_at IS NULL", clientID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r assetRepository) GetListAssets(clientID string, index, size int) ([]response.AssetResponse, error) {
	query := `
		SELECT 
			a.asset_id, a.user_client_id, a.serial_number, a.name, a.description, a.barcode,
			a.purchase_date, a.expiry_date, a.warranty_expiry_date, a.price, a.notes,
			c.asset_category_id, c.category_name, c.description AS category_description,
			s.asset_status_id, s.status_name, s.description AS status_description,
			st.stock_id, st.initial_quantity, st.latest_quantity
		FROM asset a
		JOIN asset_category c ON a.category_id = c.asset_category_id
		JOIN asset_status s ON a.status_id = s.asset_status_id
		JOIN asset_stock st ON a.asset_id = st.asset_id
		WHERE a.user_client_id = ? AND a.deleted_at IS NULL
		ORDER BY a.created_at ASC
		LIMIT ? OFFSET ?
	`

	type assetRow struct {
		AssetID             uint
		UserClientID        string
		SerialNumber        *string
		Name                string
		Description         string
		Barcode             *string
		PurchaseDateRaw     *time.Time
		ExpiryDateRaw       *time.Time
		WarrantyExpiryRaw   *time.Time
		Price               float64
		Notes               *string
		AssetStatusID       uint
		StatusName          string
		StatusDescription   string
		AssetCategoryID     uint
		CategoryName        string
		CategoryDescription string
		StockID             uint
		InitialQty          int
		LatestQty           int
	}

	var rows []assetRow
	offset := (index - 1) * size
	if err := r.db.Raw(query, clientID, size, offset).Scan(&rows).Error; err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to fetch asset list")
		return nil, err
	}

	assetResponses := make([]response.AssetResponse, len(rows))
	assetIDs := make([]uint, len(rows))
	for i, row := range rows {
		assetResponses[i] = response.AssetResponse{
			AssetID:            row.AssetID,
			UserClientID:       row.UserClientID,
			SerialNumber:       row.SerialNumber,
			Name:               row.Name,
			Description:        row.Description,
			Barcode:            row.Barcode,
			PurchaseDate:       utils.ToDateOnly(row.PurchaseDateRaw),
			ExpiryDate:         utils.ToDateOnly(row.ExpiryDateRaw),
			WarrantyExpiryDate: utils.ToDateOnly(row.WarrantyExpiryRaw),
			Price:              row.Price,
			Notes:              row.Notes,
			Status: response.AssetStatusResponse{
				AssetStatusID: row.AssetStatusID,
				StatusName:    row.StatusName,
				Description:   row.StatusDescription,
			},
			Category: response.AssetCategoryResponse{
				AssetCategoryID: row.AssetCategoryID,
				CategoryName:    row.CategoryName,
				Description:     row.CategoryDescription,
			},
			Stock: response.AssetStockResponse{
				StockID:         row.StockID,
				AssetID:         row.AssetID,
				InitialQuantity: row.InitialQty,
				LatestQuantity:  row.LatestQty,
			},
		}
		assetIDs[i] = row.AssetID
	}

	imageQuery := `
		SELECT asset_id, image_url
		FROM asset_image
		WHERE user_client_id = ? AND deleted_at IS NULL AND asset_id IN ?
	`
	var imageRows []struct {
		AssetID  uint
		ImageURL string
	}
	if err := r.db.Raw(imageQuery, clientID, assetIDs).Scan(&imageRows).Error; err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to fetch asset images")
		return nil, err
	}

	imageMap := make(map[uint][]response.AssetImageResponse)
	for _, img := range imageRows {
		imageMap[img.AssetID] = append(imageMap[img.AssetID], response.AssetImageResponse{ImageURL: img.ImageURL})
	}

	for i := range assetResponses {
		assetResponses[i].Images = imageMap[assetResponses[i].AssetID]
	}

	log.Info().Str("clientID", clientID).Int("assets_count", len(assetResponses)).Msg("✅ Successfully fetched asset list")
	return assetResponses, nil
}

func (r assetRepository) GetListAssetsByAssetGroup(clientID string, assetGroupID uint) ([]response.AssetResponse, error) {
	selectQuery := `
       SELECT 
           asset.asset_id,
           asset.user_client_id,
           asset.serial_number,
           asset.name,
           asset.description,
           asset.barcode,
           asset.purchase_date,
           asset.expiry_date,
           asset.warranty_expiry_date,
           asset.price,
           asset.notes,
           category.asset_category_id,
           category.category_name,
           category.description AS category_description,
           status.asset_status_id,
           status.status_name,
           status.description AS status_description,
           stock.stock_id,
           stock.initial_quantity,
           stock.latest_quantity           
       FROM "asset" asset
       LEFT JOIN "asset_group_asset" aga ON aga.asset_id = asset.asset_id
       INNER JOIN "asset_category" category ON asset.category_id = category.asset_category_id
       INNER JOIN "asset_status" status ON asset.status_id = status.asset_status_id
       INNER JOIN "asset_stock" stock ON asset.asset_id = stock.asset_id
       WHERE aga.asset_group_id = ? AND asset.deleted_at IS NULL
       ORDER BY asset.created_at ASC;
   `

	rows, err := r.db.Raw(selectQuery, assetGroupID).Rows()
	if err != nil {
		log.Error().Str("assetGroupID", fmt.Sprintf("%d", assetGroupID)).Err(err).Msg("❌ Failed to fetch asset list by asset group")
		return nil, err
	}

	var assetsList []response.AssetResponse
	for rows.Next() {
		var asset response.AssetResponse
		var category response.AssetCategoryResponse
		var status response.AssetStatusResponse
		var stock response.AssetStockResponse

		// Handling NULL values from SQL
		var serialNumber sql.NullString
		var barcode sql.NullString
		var description sql.NullString
		var purchaseDate sql.NullTime
		var expiryDate sql.NullTime
		var warrantyExpiryDate sql.NullTime
		var price sql.NullFloat64
		var notes sql.NullString

		err := rows.Scan(
			&asset.AssetID,
			&asset.UserClientID,
			&serialNumber,
			&asset.Name,
			&description,
			&barcode,
			&purchaseDate,
			&expiryDate,
			&warrantyExpiryDate,
			&price,
			&notes,
			&category.AssetCategoryID,
			&category.CategoryName,
			&category.Description,
			&status.AssetStatusID,
			&status.StatusName,
			&status.Description,
			&stock.StockID,
			&stock.InitialQuantity,
			&stock.LatestQuantity,
		)

		if err != nil {
			log.Error().Str("assetGroupID", fmt.Sprintf("%d", assetGroupID)).Err(err).Msg("❌ Failed to scan asset row")
			return nil, err
		}

		// Convert NULL SQL values to Go `nil`
		if serialNumber.Valid {
			asset.SerialNumber = &serialNumber.String
		}
		if barcode.Valid {
			asset.Barcode = &barcode.String
		}
		if description.Valid {
			asset.Description = description.String
		}
		if price.Valid {
			asset.Price = price.Float64
		}
		if notes.Valid {
			asset.Notes = &notes.String
		}
		if purchaseDate.Valid {
			asset.PurchaseDate = (*response.DateOnly)(&purchaseDate.Time)
		}
		if expiryDate.Valid {
			asset.ExpiryDate = (*response.DateOnly)(&expiryDate.Time)
		}
		if warrantyExpiryDate.Valid {
			asset.WarrantyExpiryDate = (*response.DateOnly)(&warrantyExpiryDate.Time)
		}

		// Assign category and status details
		asset.Category = category
		asset.Status = status
		asset.Stock = stock

		// Fetch asset images separately (handling multiple images)
		imageQuery := `
        SELECT image.image_url
        FROM "asset_image" image
        WHERE image.asset_id = ? AND user_client_id = ? AND image.deleted_at IS NULL;
    `
		imagesRows, err := r.db.Raw(imageQuery, asset.AssetID, clientID).Rows()
		if err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to fetch asset images")
			return nil, err
		}

		var images []response.AssetImageResponse
		for imagesRows.Next() {
			var img response.AssetImageResponse
			if err := imagesRows.Scan(&img.ImageURL); err != nil {
				log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset image row")
				return nil, err
			}
			images = append(images, img)
		}

		asset.Images = images

		// Append to result slice
		assetsList = append(assetsList, asset)
	}

	log.Info().Str("clientID", clientID).Int("assets_count", len(assetsList)).Msg("✅ Successfully fetched asset list")
	return assetsList, nil
}

func (r assetRepository) GetAssetResponseByID(clientID string, id uint) (*response.AssetResponse, error) {
	selectQuery := `
       SELECT 
           asset.asset_id,
           asset.user_client_id,
           asset.serial_number,
           asset.name,
           asset.description,
           asset.barcode,
           asset.purchase_date,
           asset.expiry_date,
           asset.warranty_expiry_date,
           asset.price,
           asset.notes,
           category.asset_category_id,
           category.category_name,
           category.description AS category_description,
           status.asset_status_id,
           status.status_name,
           status.description AS status_description,
           stock.stock_id,
           stock.initial_quantity,
           stock.latest_quantity           
       FROM "asset" asset
       INNER JOIN "asset_category" category ON asset.category_id = category.asset_category_id
       INNER JOIN "asset_status" status ON asset.status_id = status.asset_status_id
       INNER JOIN "asset_stock" stock ON asset.asset_id = stock.asset_id
       WHERE asset.user_client_id = ? AND asset.asset_id = ? AND asset.deleted_at IS NULL
       ORDER BY asset.created_at ASC;
   `

	rows := r.db.Raw(selectQuery, clientID, id).Row()

	var asset response.AssetResponse
	var category response.AssetCategoryResponse
	var status response.AssetStatusResponse
	var stock response.AssetStockResponse

	// Handling NULL values from SQL
	var serialNumber sql.NullString
	var barcode sql.NullString
	var description sql.NullString
	var purchaseDate sql.NullTime
	var expiryDate sql.NullTime
	var warrantyExpiryDate sql.NullTime
	var price sql.NullFloat64
	var notes sql.NullString

	err := rows.Scan(
		&asset.AssetID,
		&asset.UserClientID,
		&serialNumber,
		&asset.Name,
		&description,
		&barcode,
		&purchaseDate,
		&expiryDate,
		&warrantyExpiryDate,
		&price,
		&notes,
		&category.AssetCategoryID,
		&category.CategoryName,
		&category.Description,
		&status.AssetStatusID,
		&status.StatusName,
		&status.Description,
		&stock.StockID,
		&stock.InitialQuantity,
		&stock.LatestQuantity,
	)

	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset row")
		return nil, err
	}

	// Convert NULL SQL values to Go `nil`
	if serialNumber.Valid {
		asset.SerialNumber = &serialNumber.String
	}
	if barcode.Valid {
		asset.Barcode = &barcode.String
	}
	if description.Valid {
		asset.Description = description.String
	}
	if price.Valid {
		asset.Price = price.Float64
	}
	if notes.Valid {
		asset.Notes = &notes.String
	}
	if purchaseDate.Valid {
		asset.PurchaseDate = (*response.DateOnly)(&purchaseDate.Time)
	}
	if expiryDate.Valid {
		asset.ExpiryDate = (*response.DateOnly)(&expiryDate.Time)
	}
	if warrantyExpiryDate.Valid {
		asset.WarrantyExpiryDate = (*response.DateOnly)(&warrantyExpiryDate.Time)
	}

	// Assign category and status details
	asset.Category = category
	asset.Status = status
	asset.Stock = stock

	return &asset, nil
}

func (r assetRepository) GetAssetByID(clientID string, id uint) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("asset_id = ? AND user_client_id = ?", id, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetRepository) UpdateAsset(asset *assets.Asset, clientID string) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	// Update asset fields (only changed fields)
	if err := tx.Table(utils.TableAssetName).
		Where("asset_id = ? AND user_client_id = ?", asset.AssetID, clientID).
		Updates(asset).Error; err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetRepository) UpdateMaintenanceDateAsset(assetID uint, maintenanceDate *time.Time, clientID string) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	// Update asset fields (only changed fields)
	if err := tx.Table(utils.TableAssetName).
		Where("user_client_id = ? AND asset_id = ?", clientID, assetID).
		Updates(map[string]interface{}{"maintenance_date": maintenanceDate}).Error; err != nil {
		return fmt.Errorf("failed to update maintenance date: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetRepository) UpdateAssetStatus(assetID uint, statusID uint, clientID string) (*assets.Asset, error) {
	// Start a transaction
	tx := r.db.Begin()
	defer tx.Rollback()

	var assetOld assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("asset_id = ?", assetID).First(&assetOld).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the asset
	var asset assets.Asset
	if err := tx.Table(utils.TableAssetName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		First(&asset).Error; err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the status
	var status assets.AssetStatus
	if err := tx.Table("asset_status").Where("asset_status_id = ?", statusID).
		First(&status).Error; err != nil {
		return nil, fmt.Errorf("failed to find status: %w", err)
	}

	// Update the asset status and updated by
	if err := tx.Table(utils.TableAssetName).Model(&asset).
		Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		Updates(map[string]interface{}{
			"status_id":  statusID,
			"updated_by": clientID,
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to update asset status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &asset, nil
}

func (r assetRepository) UpdateAssetCategory(assetID uint, categoryID uint, clientID string) (*assets.Asset, error) {
	// Start a transaction
	tx := r.db.Begin()
	defer tx.Rollback()

	var assetOld assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("asset_id = ?", assetID).First(&assetOld).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the asset
	var asset assets.Asset
	if err := tx.Table(utils.TableAssetName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		First(&asset).Error; err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the category
	var category assets.AssetCategory
	if err := tx.Table("asset_category").Where("asset_category_id = ?", categoryID).
		First(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	// Update the asset category and updated by
	if err := tx.Table(utils.TableAssetName).Model(&asset).
		Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		Updates(map[string]interface{}{
			"category_id": categoryID,
			"updated_by":  clientID,
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to update asset category: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &asset, nil
}

func (r assetRepository) DeleteAsset(id uint, clientID string) error {
	if err := r.db.Table(utils.TableAssetName).Model(assets.Asset{}).
		Where("asset_id = ? AND user_client_id = ?", id, clientID).
		Updates(map[string]interface{}{"deleted_by": clientID, "deleted_at": time.Now()}).
		Delete(&assets.Asset{}).Error; err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}
	return nil
}

func (r assetRepository) GetAssetByIDForMaintenance(id uint, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("asset_id = ? AND user_client_id = ?", id, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetRepository) GetAssetByCategoryID(assetCategoryID uint, clientID string) ([]assets.Asset, error) {
	var asset []assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("category_id = ? AND user_client_id = ?", assetCategoryID, clientID).Find(&asset).Error
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (r assetRepository) GetAssetDeleted() ([]assets.Asset, error) {
	var asset []assets.Asset
	err := r.db.Unscoped().Table(utils.TableAssetName).
		Where("deleted_at IS NOT NULL").
		Find(&asset).Error
	if err != nil {
		return nil, err
	}
	return asset, nil
}
