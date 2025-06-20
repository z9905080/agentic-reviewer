package dto

import "time"

// AnalysisResponse represents the structured response for a completed CV analysis.
type AnalysisResponse struct {
	AnalysisID       string   `json:"analysis_id"`
	CVID             string   `json:"cv_id"`
	CVPath           string   `json:"cv_path"`
	JDID             string   `json:"jd_id"`
	JDPath           string   `json:"jd_path"`
	Score            float64  `json:"score"`
	Summary          string   `json:"summary"`
	MatchingKeywords []string `json:"matching_keywords"`
	MissingKeywords  []string `json:"missing_keywords"`
	AnalyzedAt       time.Time `json:"analyzed_at"`
}
