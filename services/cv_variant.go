package services

import (
	"cvbuilder/external"
	"cvbuilder/models"
	"cvbuilder/prompts"
	"cvbuilder/repos"
	"encoding/json"
	"fmt"
)

type CVVariantService struct {
	ex *external.External
	r  *repos.Repos
}

func InitCVVariantService(ex *external.External, r *repos.Repos) *CVVariantService {
	return &CVVariantService{ex: ex, r: r}
}

func (s *CVVariantService) Generate(cv *models.CV, job *models.Job, userID uint, language string) (*models.CVVariant, error) {
	cvJSON, err := json.Marshal(cv)
	if err != nil {
		return nil, fmt.Errorf("cv marshal error: %w", err)
	}

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return nil, fmt.Errorf("job marshal error: %w", err)
	}

	langLabel := language
	if language == "" || language == "auto" {
		langLabel = "the same language as the job vacancy"
	}
	langInstruction := fmt.Sprintf(
		"CRITICAL LANGUAGE RULE: You MUST write every single text field in %s. "+
			"This includes: summary, motivationLetter, every job description, AND every item in keyChanges. "+
			"Do NOT use English unless %s is English. No exceptions.",
		langLabel, langLabel,
	)

	input := langInstruction + "\n\n" + prompts.TailorCV +
		"\n\n## ORIGINAL CV\n" + string(cvJSON) +
		"\n\n## JOB VACANCY\n" + string(jobJSON)

	raw, err := s.ex.LLM.ChatCompletion(input)
	if err != nil {
		return nil, fmt.Errorf("llm error: %w", err)
	}

	var variant models.CVVariant
	if err := json.Unmarshal([]byte(stripMarkdown(raw)), &variant); err != nil {
		return nil, fmt.Errorf("json parse error: %w\nraw: %s", err, raw)
	}

	variant.UserID = userID
	variant.CVID = cv.ID
	variant.JobID = job.ID
	variant.Language = language

	if err := s.r.CVVariant.Create(&variant); err != nil {
		return nil, fmt.Errorf("save error: %w", err)
	}

	return &variant, nil
}
