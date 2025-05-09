package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetMaintenanceTypeRoutes(r *gin.Engine, middleware config.Middleware, controller assets.AssetMaintenanceTypeController) {

	routerGroup := r.Group("/v1/asset-maintenance-type")
	routerGroup.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		routerGroup.POST("/", controller.CreateMaintenanceType)
		routerGroup.GET("/:id", controller.GetMaintenanceByID)
		routerGroup.GET("/", controller.GetListMaintenanceType)
	}
}
