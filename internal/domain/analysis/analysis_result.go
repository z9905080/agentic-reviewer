package analysis

import (
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
)

// AnalysisResult represents the result of a CV to JD analysis.
type AnalysisResult struct {
	ID               string
	CV               *cv.CV
	JobDescription   *jd.JobDescription
	Score            float64 // A numerical score representing the match (e.g., 0.0 to 1.0)
	Summary          string  // A textual summary of the analysis
	MissingKeywords  []string
	MatchingKeywords []string
	// TODO: Add more detailed analysis fields if needed
}

// NewAnalysisResult creates a new AnalysisResult instance.
func NewAnalysisResult(id string, cv *cv.CV, jd *jd.JobDescription, score float64, summary string) *AnalysisResult {
	return &AnalysisResult{
		ID:             id,
		CV:             cv,
		JobDescription: jd,
		Score:          score,
		Summary:        summary,
	}
}
