package services

import "cvbuilder/external"

type WebReader struct {
	ex *external.External
}

func InitWebReader(ex *external.External) *WebReader {
	return &WebReader{
		ex,
	}
}
