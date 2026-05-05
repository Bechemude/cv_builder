package repos

import (
	"cvbuilder/db"
	"cvbuilder/models"
	"fmt"
)

type Job struct {
	db *db.DB
}

func InitJobRepo(db *db.DB) *Job {
	return &Job{db: db}
}

func (j *Job) Create(job *models.Job) error {
	if err := j.db.Postgres.Create(job).Error; err != nil {
		return fmt.Errorf("job create error: %w", err)
	}
	return nil
}

func (j *Job) GetByID(id uint) (*models.Job, error) {
	var job models.Job
	if err := j.db.Postgres.First(&job, id).Error; err != nil {
		return nil, fmt.Errorf("job get error: %w", err)
	}
	return &job, nil
}

func (j *Job) ListByUserID(userID uint) ([]models.Job, error) {
	var jobs []models.Job
	if err := j.db.Postgres.Where("user_id = ?", userID).Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("job list error: %w", err)
	}
	return jobs, nil
}
