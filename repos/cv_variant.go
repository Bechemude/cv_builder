package repos

import (
	"cvbuilder/db"
	"cvbuilder/models"
	"fmt"
)

type CVVariant struct {
	db *db.DB
}

func InitCVVariantRepo(db *db.DB) *CVVariant {
	return &CVVariant{db: db}
}

func (r *CVVariant) Create(v *models.CVVariant) error {
	if err := r.db.Postgres.Create(v).Error; err != nil {
		return fmt.Errorf("cv_variant create error: %w", err)
	}
	return nil
}

func (r *CVVariant) GetByJobID(jobID uint) (*models.CVVariant, error) {
	var v models.CVVariant
	if err := r.db.Postgres.Where("job_id = ?", jobID).First(&v).Error; err != nil {
		return nil, fmt.Errorf("cv_variant get error: %w", err)
	}
	return &v, nil
}

func (r *CVVariant) ListByUserID(userID uint) ([]models.CVVariant, error) {
	var variants []models.CVVariant
	if err := r.db.Postgres.Where("user_id = ?", userID).Find(&variants).Error; err != nil {
		return nil, fmt.Errorf("cv_variant list error: %w", err)
	}
	return variants, nil
}
