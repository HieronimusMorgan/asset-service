package assets

import (
	"asset-service/internal/controller/assets"
	"asset-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AssetWishlistRoutes(r *gin.Engine, middleware middleware.AuthMiddleware, controller assets.AssetWishlistController) {

	public := r.Group("/asset-service/v1/asset-wishlist")
	public.Use(middleware.Handler())
	{
		public.POST("/add", controller.AddWishlistAsset)
		public.POST("/update/:id", controller.UpdateWishlistAsset)
		public.GET("", controller.GetListWishlistAsset)
		public.GET("/:id", controller.GetWishlistAssetByID)
		public.DELETE("/delete/:id", controller.DeleteWishlistAsset)
	}
}
