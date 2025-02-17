package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetMaintenanceRoutes(r *gin.Engine, middleware config.Middleware, controller assets.AssetMaintenanceController) {

	routerGroup := r.Group("/asset-service/v1/asset-maintenance")
	routerGroup.Use(middleware.AuthMiddleware.Handler())
	{
		routerGroup.POST("/add-maintenance", controller.AddAssetMaintenance)
		routerGroup.GET("/:id", controller.GetMaintenanceByID)
		routerGroup.PUT("/:id", controller.UpdateMaintenance)
		routerGroup.DELETE("/:id", controller.DeleteMaintenance)
		routerGroup.GET("/asset/:asset_id", controller.GetMaintenancesByAssetID)
	}
}
