package assets

import (
	"asset-service/internal/controller/assets"
	"asset-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AssetStatusRoutes(r *gin.Engine, middleware middleware.AuthMiddleware, assetStatus assets.AssetStatusController) {

	protected := r.Group("/asset-service/v1/asset/status")
	protected.Use(middleware.Handler())
	{
		protected.GET("", assetStatus.GetListAssetStatus)
		protected.GET("/:id", assetStatus.GetAssetStatusByID)
	}

	admin := r.Group("/assets-service/v1/assets/status")
	admin.Use(middleware.Handler())
	{
		admin.POST("/add", assetStatus.AddAssetStatus)
		admin.POST("/update/:id", assetStatus.UpdateAssetStatus)
		admin.DELETE("/delete/:id", assetStatus.DeleteAssetStatus)
	}
}
