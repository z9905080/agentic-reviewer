package ports

import (
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
)

// StorageService defines the interface for storing and retrieving domain entities.
type StorageService interface {
	SaveCV(cv *cv.CV) error
	GetCV(id string) (*cv.CV, error)

	SaveJD(jd *jd.JobDescription) error
	GetJD(id string) (*jd.JobDescription, error)

	SaveAnalysisResult(result *analysis.AnalysisResult) error
	GetAnalysisResult(id string) (*analysis.AnalysisResult, error)
	GetAnalysisResultsByCV(cvID string) ([]*analysis.AnalysisResult, error)
	GetAnalysisResultsByJD(jdID string) ([]*analysis.AnalysisResult, error)
}
