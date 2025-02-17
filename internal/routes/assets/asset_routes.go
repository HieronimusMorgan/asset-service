package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetRoutes(r *gin.Engine, middleware config.Middleware, controller assets.AssetController) {

	routerGroup := r.Group("/asset-service/v1/asset")
	routerGroup.Use(middleware.AuthMiddleware.Handler())
	{
		routerGroup.POST("/add", controller.AddAsset)
		routerGroup.POST("/update/:id", controller.UpdateAsset)
		routerGroup.POST("/update-status/:id", controller.UpdateAssetStatus)
		routerGroup.POST("/update-category/:id", controller.UpdateAssetCategory)
		routerGroup.GET("", controller.GetListAsset)
		routerGroup.GET("/:id", controller.GetAssetById)
		routerGroup.DELETE("/delete/:id", controller.DeleteAsset)
	}
}
