package config

import (
	controller "asset-service/internal/controller/assets"
	"asset-service/internal/middleware"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/repository/transaction"
	users "asset-service/internal/repository/users"
	services "asset-service/internal/services/assets"
	controllercron "asset-service/internal/utils/cron/controller"
	repositorycron "asset-service/internal/utils/cron/repository"
	cron "asset-service/internal/utils/cron/service"
	"asset-service/internal/utils/jwt"
	nt "asset-service/internal/utils/nats"
	"asset-service/internal/utils/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ServerConfig holds all initialized components
type ServerConfig struct {
	Gin         *gin.Engine
	Config      *Config
	DB          *gorm.DB
	Redis       redis.RedisService
	JWTService  jwt.Service
	Cron        Cron
	Nats        Nats
	Controller  Controller
	Services    Services
	Repository  Repository
	Transaction Transaction
	Middleware  Middleware
}

// Services holds all service dependencies
type Services struct {
	AssetCategory               services.AssetCategoryService
	AssetMaintenance            services.AssetMaintenanceService
	AssetMaintenanceType        services.AssetMaintenanceTypeService
	AssetMaintenanceRecord      services.AssetMaintenanceRecordService
	Asset                       services.AssetService
	AssetStatus                 services.AssetStatusService
	AssetWishlist               services.AssetWishlistService
	AssetImage                  services.AssetImageService
	AssetGroupAssetService      services.AssetGroupAssetService
	AssetGroupMemberService     services.AssetGroupMemberService
	AssetGroupPermissionService services.AssetGroupPermissionService
	AssetGroupService           services.AssetGroupService
}

// Repository contains repository (database access objects)
type Repository struct {
	UserRepository                       users.UserRepository
	UserSettingRepository                users.UserSettingRepository
	AssetAuditLog                        repository.AssetAuditLogRepository
	AssetCategory                        repository.AssetCategoryRepository
	AssetMaintenance                     repository.AssetMaintenanceRepository
	AssetMaintenanceType                 repository.AssetMaintenanceTypeRepository
	AssetRepository                      repository.AssetRepository
	AssetStatusRepository                repository.AssetStatusRepository
	AssetWishlistRepository              repository.AssetWishlistRepository
	AssetMaintenanceRecord               repository.AssetMaintenanceRecordRepository
	AssetImageRepository                 repository.AssetImageRepository
	AssetStockRepository                 repository.AssetStockRepository
	AssetGroupRepository                 repository.AssetGroupRepository
	AssetGroupAssetRepository            repository.AssetGroupAssetRepository
	AssetGroupMemberRepository           repository.AssetGroupMemberRepository
	AssetGroupMemberPermissionRepository repository.AssetGroupMemberPermissionRepository
	AssetGroupPermissionRepository       repository.AssetGroupPermissionRepository
	AssetGroupInvitation                 repository.AssetGroupInvitationRepository
}

type Controller struct {
	AssetCategory                  controller.AssetCategoryController
	AssetMaintenance               controller.AssetMaintenanceController
	AssetMaintenanceType           controller.AssetMaintenanceTypeController
	AssetMaintenanceRecord         controller.AssetMaintenanceRecordController
	Asset                          controller.AssetController
	AssetStatus                    controller.AssetStatusController
	AssetWishlist                  controller.AssetWishlistController
	AssetGroupController           controller.AssetGroupController
	AssetGroupMemberController     controller.AssetGroupMemberController
	AssetGroupPermissionController controller.AssetGroupPermissionController
}

type Middleware struct {
	AssetMiddleware middleware.AssetMiddleware
	AdminMiddleware middleware.AdminMiddleware
}

type Cron struct {
	CronService    cron.CronService
	CronRepository repositorycron.CronRepository
	CronController controllercron.CronJobController
}

type Nats struct {
	NatsService nt.Service
}

type Transaction struct {
	AssetTransactionRepository transaction.AssetTransactionRepository
}
