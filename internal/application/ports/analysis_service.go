package ports

import (
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
)

// AnalysisService defines the interface for analyzing a CV against a Job Description.
type AnalysisService interface {
	Analyze(cv *cv.CV, jd *jd.JobDescription) (*analysis.AnalysisResult, error)
}
