package main

import (
	"asset-service/internal/database"
	"asset-service/internal/routes/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/cron/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Initialize Redis
	utils.InitializeRedis()

	// Initialize database
	db := database.InitDB()
	defer database.CloseDB(db)

	cronService := service.NewCronService(db)
	cronService.Start()
	defer cronService.Stop()

	// Setup Gin router
	r := gin.Default()

	// Register routes
	assets.AssetCategoryRoutes(r, db)
	assets.AssetStatusRoutes(r, db)
	assets.AssetRoutes(r, db)
	assets.AssetMaintenanceRoutes(r, db)

	// Run server
	log.Println("Starting server on :8081")
	err := r.Run(":8081")
	if err != nil {
		return
	}
}
