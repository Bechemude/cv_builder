package services

import (
	"cvbuilder/config"
	"cvbuilder/external"
	"cvbuilder/models"
	"cvbuilder/prompts"
	"cvbuilder/repos"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type WebReader struct {
	ex        *external.External
	r         *repos.Repos
	enricher  *Enricher
	cvVariant *CVVariantService
	c         *config.Config
}

func InitWebReader(ex *external.External, r *repos.Repos, c *config.Config) *WebReader {
	return &WebReader{
		ex:        ex,
		r:         r,
		enricher:  InitEnricher(ex, c),
		cvVariant: InitCVVariantService(ex, r, c),
	}
}

func (w *WebReader) Process(input string, userID uint, progress ProgressFunc) (*models.Job, error) {
	// [1] Fetch initial content if URL
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

	// [2] LLM mini-call: extract links + company name
	progress("🔍 Извлекаю ссылки и название компании...")
	meta, err := w.enricher.ExtractMeta(text)
	if err != nil {
		return nil, err
	}

	// [3] Parallel enrichment: fetch links + search company
	var (
		linkedPages []string
		companyInfo string
		wg          sync.WaitGroup
		mu          sync.Mutex
	)

	if len(meta.Links) > 0 {
		progress("🌐 Загружаю дополнительные страницы...")
		wg.Add(1)
		go func() {
			defer wg.Done()
			pages := w.enricher.FetchLinks(meta.Links)
			mu.Lock()
			linkedPages = pages
			mu.Unlock()
		}()
	}

	if meta.CompanyName != "" {
		progress("🏢 Ищу информацию о компании...")
		wg.Add(1)
		go func() {
			defer wg.Done()
			info := w.enricher.SearchCompany(meta.CompanyName)
			mu.Lock()
			companyInfo = info
			mu.Unlock()
		}()
	}

	wg.Wait()

	// [4] Merge all content
	merged := w.enricher.Merge(text, linkedPages, companyInfo)

	// [5] Final LLM analysis
	progress("🤖 Анализирую вакансию...")
	raw, err := w.ex.LLM.ChatCompletion(prompts.AnalyzeJob+"\n\n"+merged, w.c.ModelMain)
	if err != nil {
		return nil, fmt.Errorf("llm error: %w", err)
	}

	// [6] Save
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
