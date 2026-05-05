package services

import (
	"cvbuilder/external"
	"cvbuilder/repos"
)

type Services struct {
	User *repos.User
	CV   *repos.CV

	PDFReader       *PDFReader
	TemplateBuilder *TemplateBuilder
	WebReader       *WebReader
}

func Init(r *repos.Repos, ex *external.External) (*Services, error) {
	return &Services{
		User: r.User,
		CV:   r.CV,

		PDFReader:       InitPDFReader(ex, r),
		TemplateBuilder: InitTemplateBuilder(),
		WebReader:       InitWebReader(ex, r),
	}, nil
}
