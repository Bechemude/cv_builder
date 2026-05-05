package repos

import (
	"cvbuilder/db"
	"cvbuilder/models"
	"errors"
	"fmt"
"gorm.io/gorm"
)

type User struct {
	db *db.DB
}

func InitUserRepo(db *db.DB) *User {
	return &User{db: db}
}

func (u *User) GetByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User
	err := u.db.Postgres.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (u *User) FindOrCreate(telegramID int64, username, firstName, lastName string) (*models.User, error) {
	var user models.User

	err := u.db.Postgres.Where("telegram_id = ?", telegramID).First(&user).Error
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user find error: %w", err)
	}

	user = models.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
	}
	if err := u.db.Postgres.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("user create error: %w", err)
	}

	return &user, nil
}
