package services

import (
	"bytes"
	"cvbuilder/external"
	"cvbuilder/models"
	"cvbuilder/prompts"
	"cvbuilder/repos"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type ProgressFunc func(string)

type PDFReader struct {
	ex *external.External
	r  *repos.Repos
}

func InitPDFReader(ex *external.External, r *repos.Repos) *PDFReader {
	return &PDFReader{ex: ex, r: r}
}

func (p *PDFReader) Read(fileName string, data []byte, userID uint, progress ProgressFunc) (*models.CV, error) {
	progress("📄 Извлекаю текст из PDF...")
	text, err := extractText(data)
	if err != nil {
		return nil, fmt.Errorf("pdf extract error: %w", err)
	}
	log.Println(text)

	progress("🤖 Анализирую резюме с помощью AI...")
	raw, err := p.ex.LLM.ChatCompletion(prompts.GetInfoFromCV + "\n\n" + text)
	if err != nil {
		return nil, fmt.Errorf("llm error: %w", err)
	}

	log.Printf("llm raw response: %s", raw)

	progress("💾 Сохраняю резюме...")
	var cv models.CV
	if err := json.Unmarshal([]byte(stripMarkdown(raw)), &cv); err != nil {
		return nil, fmt.Errorf("json parse error: %w\nraw: %s", err, raw)
	}

	cv.UserID = userID

	if err := p.r.CV.Create(&cv); err != nil {
		return nil, fmt.Errorf("save error: %w", err)
	}

	return &cv, nil
}

// stripMarkdown removes ```json ... ``` or ``` ... ``` wrappers that some LLMs add.
func stripMarkdown(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSuffix(s, "```")
		s = strings.TrimSpace(s)
	}
	return s
}

func extractText(data []byte) (string, error) {
	cmd := exec.Command("pdftotext", "-", "-")
	cmd.Stdin = bytes.NewReader(data)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdftotext: %w — %s", err, strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(out.String()), nil
}
