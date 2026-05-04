package services

import "cvbuilder/repos"

type Services struct {
	r *repos.Repos

	PDFReader       *PDFReader
	TemplateBuilder *TemplateBuilder
	WebReader       *WebReader
}

func Init(r *repos.Repos) (*Services, error) {
	pdfReader := InitPDFReader()
	templateBuilder := InitTemplateBuilder()
	webReader := InitWebReader()

	return &Services{
		r: r,

		PDFReader:       pdfReader,
		TemplateBuilder: templateBuilder,
		WebReader:       webReader,
	}, nil
}
