package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetMaintenanceTypeRoutes(r *gin.Engine, middleware config.Middleware, controller assets.AssetMaintenanceTypeController) {

	routerGroup := r.Group("/asset-service/v1/asset-maintenance-type")
	routerGroup.Use(middleware.AuthMiddleware.Handler())
	{
		routerGroup.POST("/", controller.CreateMaintenanceType)
		routerGroup.GET("/:id", controller.GetMaintenanceByID)
		routerGroup.GET("/", controller.GetListMaintenanceType)
		//routerGroup.PUT("/:id", controller.UpdateMaintenance)
		//routerGroup.DELETE("/:id", controller.DeleteMaintenance)
		//routerGroup.GET("/asset/:asset_id", controller.GetMaintenancesByAssetID)
	}
}
