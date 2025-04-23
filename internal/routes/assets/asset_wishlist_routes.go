package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetWishlistRoutes(r *gin.Engine, middleware config.Middleware, controller assets.AssetWishlistController) {

	public := r.Group("/v1/asset-wishlist")
	public.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		public.POST("/add", controller.AddWishlistAsset)
		public.GET("", controller.GetListAssetWishlist)
		public.GET("/:id", controller.GetAssetWishlistByID)
		public.PUT("/:id", controller.UpdateAssetWishlist)
		public.DELETE("/:id", controller.DeleteAssetWishlist)
		public.POST("/add-to-asset/:id", controller.AddAssetWishlistToAsset)
	}
}
