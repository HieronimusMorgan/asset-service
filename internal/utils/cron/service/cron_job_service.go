package service

import (
	"asset-service/internal/services/assets"
	"asset-service/internal/utils/cron/model"
	"asset-service/internal/utils/cron/repository"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService interface {
	Start()
	Stop()
	AddCronJob(job model.CronJob)
}

// cronService implements CronService
type cronService struct {
	db                      gorm.DB
	nats                    string
	scheduler               *cron.Cron
	mu                      sync.Mutex
	jobs                    map[uint]cron.EntryID
	cronRepository          repository.CronRepository
	assetMaintenanceService assets.AssetMaintenanceService
	assetImageService       assets.AssetImageService
}

// NewCronService initializes and returns a CronService instance
func NewCronService(db gorm.DB, nats string, cronRepository repository.CronRepository, assetMaintenanceService assets.AssetMaintenanceService, image assets.AssetImageService) CronService {
	return &cronService{
		db:                      db,
		nats:                    nats,
		scheduler:               cron.New(), // Enables second-level precision
		jobs:                    make(map[uint]cron.EntryID),
		mu:                      sync.Mutex{},
		cronRepository:          cronRepository,
		assetMaintenanceService: assetMaintenanceService,
		assetImageService:       image,
	}
}

func (cs *cronService) Start() {
	cs.scheduler.Start()
	cs.loadJobsFromDB()
}

func (cs *cronService) Stop() {
	cs.scheduler.Stop()
}

func (cs *cronService) loadJobsFromDB() {
	var cronJobs []model.CronJob

	cronJobs, err := cs.cronRepository.GetCronJobs()
	if err != nil {
		log.Println("Error loading cron jobs from DB:", err)
		return
	}

	for _, job := range cronJobs {
		cs.scheduleJob(job)
	}
}

func (cs *cronService) scheduleJob(job model.CronJob) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if entryID, exists := cs.jobs[job.ID]; exists {
		cs.scheduler.Remove(entryID)
	}

	entryID, err := cs.scheduler.AddFunc(job.Schedule, func() {
		cs.executeJob(job)
	})
	if err != nil {
		log.Println("Error scheduling job:", err)
		return
	}

	cs.jobs[job.ID] = entryID
}

func (cs *cronService) executeJob(job model.CronJob) {
	now := time.Now()

	// Check for missed executions
	if !job.LastExecutedAt.IsZero() {
		expectedNextRun := job.LastExecutedAt.Add(cs.getJobInterval(job.Schedule))
		if now.After(expectedNextRun) {
			log.Printf("Job %s missed its scheduled run. Executing catch-up.\n", job.Name)
			// Handle missed execution as needed
		}
	}

	// Update the last executed time
	job.LastExecutedAt = now
	if err := cs.db.Save(&job).Error; err != nil {
		log.Println("Error updating job last executed time:", err)
	}

	// Perform the actual job task
	switch job.Name {
	case "asset_maintenance":
		err := cs.assetMaintenanceService.PerformMaintenanceCheck()
		if err != nil {
			log.Println("Error performing asset maintenance check:", err)
		}
	case "asset_image_cleanup":
		err := cs.assetImageService.Cleanup(cs.nats)
		if err != nil {
			log.Println("Error performing image cleanup:", err)
		}
	case "image_cleanup_unused":
		err := cs.assetImageService.CleanupUnusedImages(cs.nats)
		if err != nil {
			log.Println("Error performing image cleanup:", err)
		}
	default:
		log.Printf("Unknown job: %s\n", job.Name)
	}
}

func (cs *cronService) getJobInterval(schedule string) time.Duration {
	return time.Minute // Assuming a default interval of 1 minute
}

func (cs *cronService) AddCronJob(job model.CronJob) {
	if err := cs.cronRepository.Create(&job); err != nil {
		log.Println("Error creating cron job:", err)
		return
	}
}
