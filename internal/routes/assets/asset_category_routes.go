package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetCategoryRoutes(r *gin.Engine, middleware config.Middleware, controller assets.AssetCategoryController) {

	routerGroup := r.Group("/v1/asset-category")
	routerGroup.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		routerGroup.POST("/add", controller.AddAssetCategory)
		routerGroup.POST("/update/:id", controller.UpdateAssetCategory)
		routerGroup.GET("", controller.GetListAssetCategory)
		routerGroup.GET("/:id", controller.GetAssetCategoryById)
		routerGroup.DELETE("/delete/:id", controller.DeleteAssetCategory)
	}
}
