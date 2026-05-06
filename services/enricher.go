package services

import (
	"context"
	"cvbuilder/config"
	"cvbuilder/external"
	"cvbuilder/prompts"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	maxLinks     = 2
	fetchTimeout = 5 * time.Second
)

type jobMeta struct {
	CompanyName string   `json:"companyName"`
	Links       []string `json:"links"`
}

type Enricher struct {
	ex *external.External
	c  *config.Config
}

func InitEnricher(ex *external.External, c *config.Config) *Enricher {
	return &Enricher{ex: ex, c: c}
}

// ExtractMeta asks LLM to pull company name and relevant links from raw text.
func (e *Enricher) ExtractMeta(text string) (*jobMeta, error) {
	raw, err := e.ex.LLM.ChatCompletion(prompts.ExtractJobMeta+"\n\n"+text, e.c.ModelMain)
	if err != nil {
		return nil, fmt.Errorf("meta extraction error: %w", err)
	}

	var meta jobMeta
	if err := json.Unmarshal([]byte(stripMarkdown(raw)), &meta); err != nil {
		return nil, fmt.Errorf("meta parse error: %w", err)
	}

	if len(meta.Links) > maxLinks {
		meta.Links = meta.Links[:maxLinks]
	}

	return &meta, nil
}

// FetchLinks fetches pages in parallel with timeout, ignores errors.
func (e *Enricher) FetchLinks(links []string) []string {
	var (
		mu      sync.Mutex
		results []string
		wg      sync.WaitGroup
	)

	for _, link := range links {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
			defer cancel()

			ch := make(chan string, 1)
			go func() {
				text, err := e.ex.Web.FetchURL(url)
				if err == nil {
					ch <- text
				}
			}()

			select {
			case text := <-ch:
				mu.Lock()
				results = append(results, fmt.Sprintf("=== %s ===\n%s", url, text))
				mu.Unlock()
			case <-ctx.Done():
				// timeout — skip
			}
		}(link)
	}

	wg.Wait()
	return results
}

// SearchCompany fetches company info from the web, ignores errors.
func (e *Enricher) SearchCompany(name string) string {
	if name == "" {
		return ""
	}
	text, err := e.ex.Web.SearchCompany(name)
	if err != nil {
		return ""
	}
	return "=== Company info: " + name + " ===\n" + text
}

// Merge combines all sources into one context string for the final LLM call.
func (e *Enricher) Merge(original string, linkedPages []string, companyInfo string) string {
	var sb strings.Builder

	sb.WriteString("=== Job Vacancy ===\n")
	sb.WriteString(original)

	for _, page := range linkedPages {
		sb.WriteString("\n\n")
		sb.WriteString(page)
	}

	if companyInfo != "" {
		sb.WriteString("\n\n")
		sb.WriteString(companyInfo)
	}

	return sb.String()
}
