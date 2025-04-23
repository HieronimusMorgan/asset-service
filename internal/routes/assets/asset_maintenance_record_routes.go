package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetMaintenanceRecordRoutes(r *gin.Engine, middleware config.Middleware, controller assets.AssetMaintenanceRecordController) {

	routerGroup := r.Group("/v1/asset-maintenance-record")
	routerGroup.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		routerGroup.GET("/", controller.GetMaintenanceRecord)
	}
}
