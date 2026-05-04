package services

import (
	"cvbuilder/external"
	"cvbuilder/repos"
)

type Services struct {
	r *repos.Repos

	PDFReader       *PDFReader
	TemplateBuilder *TemplateBuilder
	WebReader       *WebReader
}

func Init(r *repos.Repos, ex *external.External) (*Services, error) {
	pdfReader := InitPDFReader()
	templateBuilder := InitTemplateBuilder()
	webReader := InitWebReader(ex)

	return &Services{
		r: r,

		PDFReader:       pdfReader,
		TemplateBuilder: templateBuilder,
		WebReader:       webReader,
	}, nil
}
