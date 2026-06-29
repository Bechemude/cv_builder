package services

import (
	"cvbuilder/config"
	"cvbuilder/external"
	"cvbuilder/models"
	"cvbuilder/prompts"
	"cvbuilder/repos"
	"encoding/json"
	"fmt"
	"log"
)

type CVVariantService struct {
	ex *external.External
	r  *repos.Repos
	c  *config.Config
}

func InitCVVariantService(ex *external.External, r *repos.Repos, c *config.Config) *CVVariantService {
	return &CVVariantService{ex: ex, r: r, c: c}
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

	var variant models.CVVariant
	var lastErr error

	for attempt := range 3 {
		raw, err := s.ex.LLM.ChatCompletion(input, s.c.ModelMain)
		if err != nil {
			lastErr = fmt.Errorf("llm error: %w", err)
			continue
		}

		cleaned := cleanLLMResponse(raw)

		if err := json.Unmarshal([]byte(cleaned), &variant); err != nil {
			log.Printf("json parse error (attempt %d/3): %v\nraw: %s\ncleaned: %s",
				attempt+1, err, raw, cleaned)
			lastErr = fmt.Errorf("json parse error: %w", err)
			continue
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

	return nil, fmt.Errorf("generate variant error: %w", lastErr)
}
