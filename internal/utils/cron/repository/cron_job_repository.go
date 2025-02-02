package repository

import (
	"asset-service/internal/utils/cron/model"
	"gorm.io/gorm"
)

type CronRepository struct {
	db *gorm.DB
}

func NewCronRepository(db *gorm.DB) *CronRepository {
	return &CronRepository{db: db}
}

func (r *CronRepository) GetCronJobs() ([]model.CronJob, error) {
	var cronJobs []model.CronJob
	err := r.db.Find(&cronJobs).Error
	if err != nil {
		return nil, err
	}
	return cronJobs, nil
}

func (r *CronRepository) GetCronJobByID(id uint) (model.CronJob, error) {
	var cronJob model.CronJob
	err := r.db.Where("id = ?", id).First(&cronJob).Error
	if err != nil {
		return model.CronJob{}, err
	}
	return cronJob, nil
}

func (r *CronRepository) CreateCronJob(cronJob *model.CronJob) error {
	err := r.db.Create(&cronJob).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CronRepository) UpdateCronJob(cronJob *model.CronJob) error {
	err := r.db.Table("my-home.cron_job").Save(&cronJob).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CronRepository) DeleteCronJob(id uint) error {
	err := r.db.Table("my-home.cron_job").Delete(&model.CronJob{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CronRepository) GetCronJobByJobName(jobName string) (model.CronJob, error) {
	var cronJob model.CronJob
	err := r.db.Table("my-home.cron_job").Where("job_name = ?", jobName).First(&cronJob).Error
	if err != nil {
		return model.CronJob{}, err
	}
	return cronJob, nil
}

func (r *CronRepository) deleteCronJobByID(id uint) error {
	err := r.db.Table("my-home.cron_job").Delete(&model.CronJob{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
