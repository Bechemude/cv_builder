package services

import (
	"cvbuilder/external"
	"cvbuilder/models"
	"cvbuilder/prompts"
	"cvbuilder/repos"
	"encoding/json"
	"fmt"
	"strings"
)

type WebReader struct {
	ex *external.External
	r  *repos.Repos
}

func InitWebReader(ex *external.External, r *repos.Repos) *WebReader {
	return &WebReader{ex: ex, r: r}
}

func (w *WebReader) Process(input string, userID uint, progress ProgressFunc) (*models.Job, error) {
	var text string

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		progress("🌐 Загружаю страницу вакансии...")
		raw, err := w.ex.Web.FetchURL(input)
		if err != nil {
			return nil, fmt.Errorf("fetch error: %w", err)
		}
		text = raw
	} else {
		text = input
	}

	progress("🤖 Анализирую вакансию...")
	raw, err := w.ex.LLM.ChatCompletion(prompts.AnalyzeJob + "\n\n" + text)
	if err != nil {
		return nil, fmt.Errorf("llm error: %w", err)
	}

	progress("💾 Сохраняю вакансию...")
	var job models.Job
	if err := json.Unmarshal([]byte(stripMarkdown(raw)), &job); err != nil {
		return nil, fmt.Errorf("json parse error: %w\nraw: %s", err, raw)
	}

	job.UserID = userID
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		job.SourceUrl = input
	}

	if err := w.r.Job.Create(&job); err != nil {
		return nil, fmt.Errorf("save error: %w", err)
	}

	return &job, nil
}
