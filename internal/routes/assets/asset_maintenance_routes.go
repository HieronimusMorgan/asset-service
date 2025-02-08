package assets

import (
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssetMaintenanceRoutes(r *gin.Engine, db *gorm.DB) {
	assetHandler := assets.NewAssetMaintenanceController(db)

	assetMaintenanceRoutes := r.Group("/asset-service/v1/asset-maintenance")
	{
		assetMaintenanceRoutes.POST("/", assetHandler.CreateMaintenance)
		assetMaintenanceRoutes.GET("/:id", assetHandler.GetMaintenanceByID)
		assetMaintenanceRoutes.PUT("/:id", assetHandler.UpdateMaintenance)
		assetMaintenanceRoutes.DELETE("/:id", assetHandler.DeleteMaintenance)
		assetMaintenanceRoutes.GET("/asset/:asset_id", assetHandler.GetMaintenancesByAssetID)
	}
}
