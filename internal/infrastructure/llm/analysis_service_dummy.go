package llm

import (
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"fmt"
	"github.com/google/uuid"
)

// DummyAnalysisService is a placeholder implementation for AnalysisService.
type DummyAnalysisService struct{}

// NewDummyAnalysisService creates a new DummyAnalysisService.
func NewDummyAnalysisService() *DummyAnalysisService {
	return &DummyAnalysisService{}
}

// Analyze simulates analyzing a CV against a Job Description.
func (s *DummyAnalysisService) Analyze(cv *cv.CV, jd *jd.JobDescription) (*analysis.AnalysisResult, error) {
	if cv == nil || jd == nil {
		return nil, fmt.Errorf("CV and Job Description cannot be nil")
	}
	fmt.Printf("Simulating analysis for CV ID: %s and JD ID: %s\n", cv.ID, jd.ID)

	// Simulate analysis result
	score := 0.75 // Example score
	summary := fmt.Sprintf("This is a dummy analysis summary for CV %s and JD %s. The candidate seems like a good fit.", cv.ID, jd.ID)

	return &analysis.AnalysisResult{
		ID:               uuid.NewString(),
		CV:               cv,
		JobDescription:   jd,
		Score:            score,
		Summary:          summary,
		MissingKeywords:  []string{"dummy_missing_keyword_1", "dummy_missing_keyword_2"},
		MatchingKeywords: []string{"dummy_matching_keyword_1", "dummy_matching_keyword_2", "dummy_matching_keyword_3"},
	}, nil
}
