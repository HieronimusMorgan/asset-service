package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetMaintenanceRoutes(r *gin.Engine, middleware config.Middleware, controller assets.AssetMaintenanceController) {

	routerGroup := r.Group("/v1/asset-maintenance")
	routerGroup.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		routerGroup.POST("/add-maintenance", controller.AddAssetMaintenance)
		routerGroup.POST("/perform-maintenance", controller.PerformMaintenance)
		routerGroup.GET("/:id", controller.GetMaintenanceByID)
		routerGroup.PUT("/:id", controller.UpdateMaintenance)
		routerGroup.DELETE("/:id", controller.DeleteMaintenance)
		routerGroup.GET("/asset/:asset_id", controller.GetMaintenancesByAssetID)
	}
}
