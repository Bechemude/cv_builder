package models

import "gorm.io/gorm"

type CVVariant struct {
	gorm.Model

	UserID uint `json:"-" gorm:"column:user_id"`
	CVID   uint `json:"-" gorm:"column:cv_id"`
	JobID  uint `json:"-" gorm:"column:job_id"`

	// Tailored narrative fields
	Summary          string       `json:"summary"          gorm:"column:summary"`
	MotivationLetter string       `json:"motivationLetter" gorm:"column:motivation_letter"`
	JobsHistory      []JobHistory `json:"jobsHistory"      gorm:"serializer:json"`
	Tags             []string     `json:"tags"             gorm:"serializer:json"`

	// Analysis
	MatchScore int      `json:"matchScore" gorm:"column:match_score"`
	KeyChanges []string `json:"keyChanges" gorm:"serializer:json"`
	Language   string   `json:"language"   gorm:"column:language"`
}

// BuildCV assembles a complete CV from the original (non-narrative fields)
// and this variant (tailored narrative fields).
func (v *CVVariant) BuildCV(original *CV) *CV {
	return &CV{
		UserID:           original.UserID,
		FirstName:        original.FirstName,
		LastName:         original.LastName,
		DOB:              original.DOB,
		Summary:          v.Summary,
		MotivationLetter: v.MotivationLetter,
		JobsHistory:      v.JobsHistory,
		Tags:             v.Tags,
	}
}
