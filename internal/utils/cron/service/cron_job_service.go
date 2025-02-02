package service

import (
	"asset-service/internal/utils/cron/model"
	"asset-service/internal/utils/cron/repository"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService struct {
	db             *gorm.DB
	scheduler      *cron.Cron
	jobs           map[uint]cron.EntryID
	mu             sync.Mutex
	cronRepository *repository.CronRepository
}

func NewCronService(db *gorm.DB) *CronService {
	// Configure the cron scheduler to accept cron expressions with seconds
	scheduler := cron.New()

	return &CronService{
		db:             db,
		scheduler:      scheduler,
		jobs:           make(map[uint]cron.EntryID),
		cronRepository: repository.NewCronRepository(db),
	}
}

func (cs *CronService) Start() {
	cs.scheduler.Start()
	cs.loadJobsFromDB()
}

func (cs *CronService) Stop() {
	cs.scheduler.Stop()
}

func (cs *CronService) loadJobsFromDB() {
	var cronJobs []model.CronJob
	if err := cs.db.Where("is_active = ?", true).Find(&cronJobs).Error; err != nil {
		log.Println("Error loading cron jobs from DB:", err)
		return
	}

	for _, job := range cronJobs {
		cs.scheduleJob(job)
	}
}

func (cs *CronService) scheduleJob(job model.CronJob) {
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

func (cs *CronService) executeJob(job model.CronJob) {
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
	log.Printf("Executing job: %s\n", job.Name)
	// Add your job logic here
}

func (cs *CronService) getJobInterval(schedule string) time.Duration {
	// Parse the cron schedule to determine the interval
	// This is a simplified example; you'll need to implement parsing based on your cron library
	return time.Minute // Assuming a default interval of 1 minute
}
