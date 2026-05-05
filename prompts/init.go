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

//go:embed EXTRACT_JOB_META.md
var ExtractJobMeta string

//go:embed TAILOR_CV.md
var tailorCVTemplate string

var GetInfoFromCV string
var AnalyzeJob string
var TailorCV string

func init() {
	GetInfoFromCV = strings.ReplaceAll(getInfoFromCVTemplate, "{{SCHEMA}}", models.CVSchema())
	AnalyzeJob = strings.ReplaceAll(analyzeJobTemplate, "{{SCHEMA}}", models.JobSchema())
	TailorCV = strings.ReplaceAll(tailorCVTemplate, "{{SCHEMA}}", models.CVVariantSchema())
}
