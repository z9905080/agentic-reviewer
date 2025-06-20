package mocks

import (
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"cv-analyzer/internal/application/ports" // For AnalysisService interface
	"github.com/stretchr/testify/mock"
)

type MockAnalysisService struct {
	mock.Mock
}

func (m *MockAnalysisService) Analyze(cvData *cv.CV, jdData *jd.JobDescription) (*analysis.AnalysisResult, error) {
	args := m.Called(cvData, jdData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*analysis.AnalysisResult), args.Error(1)
}
var _ ports.AnalysisService = (*MockAnalysisService)(nil)
