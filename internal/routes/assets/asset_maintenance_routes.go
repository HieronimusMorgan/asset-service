package assets

import (
	"asset-service/internal/controller/assets"
	"asset-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AssetMaintenanceRoutes(r *gin.Engine, middleware middleware.AuthMiddleware, controller assets.AssetMaintenanceController) {

	assetMaintenanceRoutes := r.Group("/asset-service/v1/asset-maintenance")
	assetMaintenanceRoutes.Use(middleware.Handler())
	{
		assetMaintenanceRoutes.POST("/", controller.CreateMaintenance)
		assetMaintenanceRoutes.GET("/:id", controller.GetMaintenanceByID)
		assetMaintenanceRoutes.PUT("/:id", controller.UpdateMaintenance)
		assetMaintenanceRoutes.DELETE("/:id", controller.DeleteMaintenance)
		assetMaintenanceRoutes.GET("/asset/:asset_id", controller.GetMaintenancesByAssetID)
	}
}
