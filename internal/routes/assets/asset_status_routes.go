package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetStatusRoutes(r *gin.Engine, middleware config.Middleware, assetStatus assets.AssetStatusController) {

	routerGroup := r.Group("/asset-service/v1/asset/status")
	routerGroup.Use(middleware.AuthMiddleware.Handler())
	{
		routerGroup.GET("", assetStatus.GetListAssetStatus)
		routerGroup.GET("/:id", assetStatus.GetAssetStatusByID)
	}

	admin := r.Group("/assets-service/v1/assets/status")
	admin.Use(middleware.AdminMiddleware.Handler())
	{
		admin.POST("/add", assetStatus.AddAssetStatus)
		admin.POST("/update/:id", assetStatus.UpdateAssetStatus)
		admin.DELETE("/delete/:id", assetStatus.DeleteAssetStatus)
	}
}
