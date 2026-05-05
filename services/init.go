package services

import (
	"cvbuilder/external"
	"cvbuilder/repos"
	_ "embed"
	"fmt"
)

//go:embed templates/cv.html
var cvTemplate string

type Services struct {
	User      *repos.User
	CV        *repos.CV
	Job       *repos.Job
	CVVariant *CVVariantService

	PDFReader       *PDFReader
	PDFGenerator    *PDFGenerator
	TemplateBuilder *TemplateBuilder
	WebReader       *WebReader
}

func Init(r *repos.Repos, ex *external.External) (*Services, error) {
	pdfGen, err := InitPDFGenerator(cvTemplate)
	if err != nil {
		return nil, fmt.Errorf("pdf generator init error: %w", err)
	}

	return &Services{
		User:      r.User,
		CV:        r.CV,
		Job:       r.Job,
		CVVariant: InitCVVariantService(ex, r),

		PDFReader:       InitPDFReader(ex, r),
		PDFGenerator:    pdfGen,
		TemplateBuilder: InitTemplateBuilder(),
		WebReader:       InitWebReader(ex, r),
	}, nil
}
