package assets

import (
	"asset-service/internal/controller/assets"
	"asset-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AssetRoutes(r *gin.Engine, middleware middleware.AuthMiddleware, controller assets.AssetController) {

	public := r.Group("/asset-service/v1/asset")
	public.Use(middleware.Handler())
	{
		public.POST("/add", controller.AddAsset)
		//public.POST("/wishlist", controller.AddWishlistAsset)
		public.POST("/update/:id", controller.UpdateAsset)
		public.POST("/update-status/:id", controller.UpdateAssetStatus)
		public.POST("/update-category/:id", controller.UpdateAssetCategory)
		public.GET("", controller.GetListAsset)
		public.GET("/:id", controller.GetAssetById)
		public.DELETE("/delete/:id", controller.DeleteAsset)
	}
}
