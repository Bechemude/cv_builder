package prompts

import (
	_ "embed"
	"strings"

	"cvbuilder/models"
)

//go:embed GET_INFO_FROM_CV.md
var getInfoFromCVTemplate string

//go:embed ANALYZE_JOB.md
var analyzeJobTemplate string

var GetInfoFromCV string
var AnalyzeJob string

func init() {
	GetInfoFromCV = strings.ReplaceAll(getInfoFromCVTemplate, "{{SCHEMA}}", models.CVSchema())
	AnalyzeJob = strings.ReplaceAll(analyzeJobTemplate, "{{SCHEMA}}", models.JobSchema())
}
