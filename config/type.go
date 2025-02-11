package config

import (
	controller "asset-service/internal/controller/assets"
	"asset-service/internal/middleware"
	repository "asset-service/internal/repository/assets"
	services "asset-service/internal/services/assets"
	"asset-service/internal/utils"
	controllercron "asset-service/internal/utils/cron/controller"
	repositorycron "asset-service/internal/utils/cron/repository"
	cron "asset-service/internal/utils/cron/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ServerConfig holds all initialized components
type ServerConfig struct {
	Gin        *gin.Engine
	Config     *Config
	DB         *gorm.DB
	Redis      utils.RedisService
	JWTService utils.JWTService
	Cron       Cron
	Controller Controller
	Services   Services
	Repository Repository
	Middleware Middleware
}

// Services holds all service dependencies
type Services struct {
	AssetCategory        services.AssetCategoryService
	AssetMaintenance     services.AssetMaintenanceService
	AssetMaintenanceType services.AssetMaintenanceTypeService
	Asset                services.AssetService
	AssetStatus          services.AssetStatusService
	AssetWishlist        services.AssetWishlistService
}

// Repository contains repository (database access objects)
type Repository struct {
	AssetAuditLog           repository.AssetAuditLogRepository
	AssetCategory           repository.AssetCategoryRepository
	AssetMaintenance        repository.AssetMaintenanceRepository
	AssetMaintenanceType    repository.AssetMaintenanceTypeRepository
	AssetRepository         repository.AssetRepository
	AssetStatusRepository   repository.AssetStatusRepository
	AssetWishlistRepository repository.AssetWishlistRepository
}

type Controller struct {
	AssetCategory        controller.AssetCategoryController
	AssetMaintenance     controller.AssetMaintenanceController
	AssetMaintenanceType controller.AssetMaintenanceTypeController
	Asset                controller.AssetController
	AssetStatus          controller.AssetStatusController
	AssetWishlist        controller.AssetWishlistController
}

type Middleware struct {
	AuthMiddleware  middleware.AuthMiddleware
	AdminMiddleware middleware.AdminMiddleware
}

type Cron struct {
	CronService    cron.CronService
	CronRepository repositorycron.CronRepository
	CronController controllercron.CronJobController
}
