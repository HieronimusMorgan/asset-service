package controller

import (
	"asset-service/internal/utils/cron/service"
	"gorm.io/gorm"
)

type CronJobController struct {
	cronJobService *service.CronService
}

func NewCronJobController(db *gorm.DB) *CronJobController {
	s := service.NewCronService(db)
	return &CronJobController{cronJobService: s}
}

func (h CronJobController) AddCronJob() {
	//h.cronJobService.AddCronJob()
}
