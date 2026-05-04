package handlers

import "cvbuilder/services"

type Handlers struct {
	s *services.Services
}

func Init(s *services.Services) (*Handlers, error) {
	return &Handlers{
		s,
	}, nil
}
