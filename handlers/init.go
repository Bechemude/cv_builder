package handlers

import (
	"cvbuilder/external"
	"cvbuilder/services"
)

type Handlers struct {
	s  *services.Services
	ex *external.External
}

func Init(s *services.Services, ex *external.External) (*Handlers, error) {
	return &Handlers{
		s,
		ex,
	}, nil
}
