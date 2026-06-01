package services

import (
	"bytes"
	"context"
	"cvbuilder/models"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//go:generate echo "templates are embedded via init"

type PDFGenerator struct {
	tmpl *template.Template
}

// cvLabels holds localized section titles.
type cvLabels struct {
	Summary    string
	Experience string
	Skills     string
	Motivation string
	Present    string
}

// cvJobData is the per-job view model used in the template.
type cvJobData struct {
	Title       string
	Position    string
	CompanyName string
	Description string
	Tags        []string
	StartStr    string
	EndStr      string
}

// cvData is the full view model passed to the HTML template.
type cvData struct {
	FirstName        string
	LastName         string
	CurrentPosition  string
	Summary          string
	MotivationLetter string
	Skills           []string
	Jobs             []cvJobData
	Language         string
	Labels           cvLabels
}

func InitPDFGenerator(tmplContent string) (*PDFGenerator, error) {
	tmpl, err := template.New("cv").Parse(tmplContent)
	if err != nil {
		return nil, fmt.Errorf("pdf template parse error: %w", err)
	}
	return &PDFGenerator{tmpl: tmpl}, nil
}

// Render generates a PDF from a CV model and returns the raw bytes.
func (g *PDFGenerator) Render(cv *models.CV, language string) ([]byte, error) {
	data := buildCVData(cv, language)

	var buf bytes.Buffer
	if err := g.tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("template execute error: %w", err)
	}
	htmlStr := buf.String()

	return renderToPDF(htmlStr)
}

// buildCVData converts a models.CV into the template view model.
func buildCVData(cv *models.CV, language string) cvData {
	labels := labelsFor(language)

	// Build jobs
	jobs := make([]cvJobData, 0, len(cv.JobsHistory))
	for _, j := range cv.JobsHistory {
		jobs = append(jobs, cvJobData{
			Title:       j.Title,
			Position:    j.Position,
			CompanyName: j.CompanyName,
			Description: j.Description,
			Tags:        j.Tags,
			StartStr:    formatFlexTime(j.Start),
			EndStr:      formatFlexTimeEnd(j.End, labels.Present),
		})
	}

	// Current position = most recent job's position
	currentPos := ""
	if len(jobs) > 0 {
		currentPos = jobs[0].Position
		if jobs[0].Title != "" && jobs[0].Title != currentPos {
			currentPos = jobs[0].Title
		}
	}

	return cvData{
		FirstName:        cv.FirstName,
		LastName:         cv.LastName,
		CurrentPosition:  currentPos,
		Summary:          cv.Summary,
		MotivationLetter: cv.MotivationLetter,
		Skills:           cv.Tags,
		Jobs:             jobs,
		Language:         language,
		Labels:           labels,
	}
}

// labelsFor returns localized section titles.
func labelsFor(language string) cvLabels {
	switch language {
	case "ru":
		return cvLabels{
			Summary:    "О себе",
			Experience: "Опыт работы",
			Skills:     "Навыки",
			Motivation: "Сопроводительное письмо",
			Present:    "н.в.",
		}
	default:
		return cvLabels{
			Summary:    "Summary",
			Experience: "Experience",
			Skills:     "Skills",
			Motivation: "Cover Letter",
			Present:    "Present",
		}
	}
}

// formatFlexTime formats a FlexTime as "Jan 2006".
func formatFlexTime(ft models.FlexTime) string {
	if ft.T == nil {
		return ""
	}
	return ft.T.Format("Jan 2006")
}

// formatFlexTimeEnd formats end date; returns presentLabel if nil (current job).
func formatFlexTimeEnd(ft models.FlexTime, presentLabel string) string {
	if ft.T == nil {
		return presentLabel
	}
	return ft.T.Format("Jan 2006")
}

// renderToPDF takes an HTML string, writes it to a temp file, and renders via chromedp.
func renderToPDF(htmlStr string) ([]byte, error) {
	// Write to temp file so chromedp can load it as file://
	tmpFile, err := os.CreateTemp("", "cv-*.html")
	if err != nil {
		return nil, fmt.Errorf("temp file error: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(htmlStr); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("temp file write error: %w", err)
	}
	tmpFile.Close()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var pdfBuf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate("file://"+tmpFile.Name()),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).   // A4 width in inches
				WithPaperHeight(11.69). // A4 height in inches
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("chromedp render error: %w", err)
	}

	return pdfBuf, nil
}
