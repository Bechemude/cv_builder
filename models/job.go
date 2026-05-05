package models

import "gorm.io/gorm"

type Job struct {
	gorm.Model

	UserID uint `json:"-" gorm:"column:user_id"`

	// Basic info
	Title       string `json:"title"       gorm:"column:title"`
	Position    string `json:"position"    gorm:"column:position"`
	CompanyName string `json:"companyName" gorm:"column:company_name"`
	CompanyUrl  string `json:"companyUrl"  gorm:"column:company_url"`
	SourceUrl   string `json:"sourceUrl"   gorm:"column:source_url"`

	// Role details
	Description     string `json:"description"     gorm:"column:description"`
	Responsibilities string `json:"responsibilities" gorm:"column:responsibilities"`
	Seniority       string `json:"seniority"       gorm:"column:seniority"`
	EmploymentType  string `json:"employmentType"  gorm:"column:employment_type"`

	// Location
	Location  string `json:"location"  gorm:"column:location"`
	Remote    string `json:"remote"    gorm:"column:remote"`

	// Compensation
	SalaryFrom int    `json:"salaryFrom" gorm:"column:salary_from"`
	SalaryTo   int    `json:"salaryTo"   gorm:"column:salary_to"`
	Currency   string `json:"currency"   gorm:"column:currency"`

	// Skills & stack
	SkillsRequired   []string `json:"skillsRequired"   gorm:"serializer:json"`
	SkillsNiceToHave []string `json:"skillsNiceToHave" gorm:"serializer:json"`
	TechStack        []string `json:"techStack"        gorm:"serializer:json"`

	// Company
	CompanySize   string `json:"companySize"   gorm:"column:company_size"`
	CompanyStage  string `json:"companyStage"  gorm:"column:company_stage"`
	Industry      string `json:"industry"      gorm:"column:industry"`
	TeamSize      string `json:"teamSize"      gorm:"column:team_size"`
	Benefits      []string `json:"benefits"    gorm:"serializer:json"`

	// Analysis
	UniqueTraits []string `json:"uniqueTraits" gorm:"serializer:json"`
	RedFlags     []string `json:"redFlags"     gorm:"serializer:json"`
	Language     string   `json:"language"     gorm:"column:language"`
}
