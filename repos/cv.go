package repos

import (
	"cvbuilder/db"
	"cvbuilder/models"
	"fmt"
)

type CV struct {
	db *db.DB
}

func InitCVRepo(db *db.DB) *CV {
	return &CV{db: db}
}

func (c *CV) Create(cv *models.CV) error {
	result := c.db.Postgres.Create(cv)
	if result.Error != nil {
		return fmt.Errorf("cv create error: %w", result.Error)
	}
	return nil
}

func (c *CV) GetByID(id uint) (*models.CV, error) {
	var cv models.CV
	result := c.db.Postgres.Preload("JobsHistory").First(&cv, id)
	if result.Error != nil {
		return nil, fmt.Errorf("cv get error: %w", result.Error)
	}
	return &cv, nil
}

func (c *CV) Update(cv *models.CV) error {
	result := c.db.Postgres.Save(cv)
	if result.Error != nil {
		return fmt.Errorf("cv update error: %w", result.Error)
	}
	return nil
}

func (c *CV) LoadJobsHistory(cv *models.CV) error {
	return c.db.Postgres.Preload("JobsHistory").First(cv, cv.ID).Error
}

func (c *CV) ListByUserID(userID uint) ([]models.CV, error) {
	var cvs []models.CV
	if err := c.db.Postgres.Where("user_id = ?", userID).Find(&cvs).Error; err != nil {
		return nil, fmt.Errorf("cv list error: %w", err)
	}
	return cvs, nil
}

func (c *CV) Delete(id uint) error {
	result := c.db.Postgres.Delete(&models.CV{}, id)
	if result.Error != nil {
		return fmt.Errorf("cv delete error: %w", result.Error)
	}
	return nil
}
