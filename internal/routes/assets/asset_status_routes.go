package assets

import (
	"asset-service/internal/controller/assets"
	"asset-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssetStatusRoutes(r *gin.Engine, db *gorm.DB) {
	assetStatusHandler := assets.NewAssetStatusController(db)

	protected := r.Group("/asset-service/v1/asset/status")
	protected.Use(middleware.Middleware())
	{
		protected.GET("", assetStatusHandler.GetListAssetStatus)
		protected.GET("/:id", assetStatusHandler.GetAssetStatusByID)
	}

	admin := r.Group("/assets-service/v1/assets/status")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.POST("/add", assetStatusHandler.AddAssetStatus)
		admin.POST("/update/:id", assetStatusHandler.UpdateAssetStatus)
		admin.DELETE("/delete/:id", assetStatusHandler.DeleteAssetStatus)
	}
}
