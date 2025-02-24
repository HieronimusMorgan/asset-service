package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetWishlistRoutes(r *gin.Engine, middleware config.Middleware, controller assets.AssetWishlistController) {

	public := r.Group("/v1/asset-wishlist")
	public.Use(middleware.AuthMiddleware.Handler())
	{
		public.POST("/add", controller.AddWishlistAsset)
		public.POST("/update/:id", controller.UpdateWishlistAsset)
		public.GET("", controller.GetListWishlistAsset)
		public.GET("/:id", controller.GetWishlistAssetByID)
		public.DELETE("/delete/:id", controller.DeleteWishlistAsset)
	}
}
