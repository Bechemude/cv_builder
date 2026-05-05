package models

import "gorm.io/gorm"

type CV struct {
	gorm.Model

	UserID           uint         `json:"-"                gorm:"column:user_id"`
	FirstName        string       `json:"firstName"        gorm:"column:first_name"`
	LastName         string       `json:"lastName"         gorm:"column:last_name"`
	DOB              FlexTime     `json:"dob"              gorm:"column:dob"`
	MotivationLetter string       `json:"motivationLetter" gorm:"column:motivation_letter"`
	Summary          string       `json:"summary"          gorm:"column:summary"`
	JobsHistory      []JobHistory `json:"jobsHistory"      gorm:"foreignKey:CVID"`
	Tags             []string     `json:"tags"             gorm:"serializer:json"`
}

type JobHistory struct {
	gorm.Model

	CVID        uint       `json:"-"            gorm:"column:cv_id"`
	Title       string     `json:"title"        gorm:"column:title"`
	Description string     `json:"description"  gorm:"column:description"`
	CompanyName string     `json:"companyName"  gorm:"column:company_name"`
	CompanyUrl  string     `json:"companyUrl"   gorm:"column:company_url"`
	Position    string     `json:"position"     gorm:"column:position"`
	Start       FlexTime   `json:"start"        gorm:"column:start"`
	End         FlexTime   `json:"end"          gorm:"column:end"`
	Tags        []string   `json:"tags"         gorm:"serializer:json"`
}
