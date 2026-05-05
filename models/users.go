package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	TelegramID int64  `json:"telegramId" gorm:"column:telegram_id;uniqueIndex"`
	Username   string `json:"username"   gorm:"column:username"`
	FirstName  string `json:"firstName"  gorm:"column:first_name"`
	LastName   string `json:"lastName"   gorm:"column:last_name"`

	CVs []CV `json:"cvs" gorm:"foreignKey:UserID"`
}
