package dto

// AnalysisRequest represents the request payload for initiating a CV analysis.
type AnalysisRequest struct {
	CVPath string `json:"cv_path" binding:"required"`
	JDPath string `json:"jd_path" binding:"required"`
}
