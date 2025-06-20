package mocks

import (
	"cv-analyzer/internal/application/ports" // For StorageService interface
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"github.com/stretchr/testify/mock"
)

type MockStorageService struct {
	mock.Mock
}

// SaveCV mocks the SaveCV method
func (m *MockStorageService) SaveCV(cvData *cv.CV) error {
	args := m.Called(cvData)
	return args.Error(0)
}

// GetCV mocks the GetCV method
func (m *MockStorageService) GetCV(id string) (*cv.CV, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cv.CV), args.Error(1)
}

// SaveJD mocks the SaveJD method
func (m *MockStorageService) SaveJD(jdData *jd.JobDescription) error {
	args := m.Called(jdData)
	return args.Error(0)
}

// GetJD mocks the GetJD method
func (m *MockStorageService) GetJD(id string) (*jd.JobDescription, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jd.JobDescription), args.Error(1)
}

// SaveAnalysisResult mocks the SaveAnalysisResult method
func (m *MockStorageService) SaveAnalysisResult(result *analysis.AnalysisResult) error {
	args := m.Called(result)
	return args.Error(0)
}

// GetAnalysisResult mocks the GetAnalysisResult method
func (m *MockStorageService) GetAnalysisResult(id string) (*analysis.AnalysisResult, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*analysis.AnalysisResult), args.Error(1)
}

// GetAnalysisResultsByCV mocks the GetAnalysisResultsByCV method
func (m *MockStorageService) GetAnalysisResultsByCV(cvID string) ([]*analysis.AnalysisResult, error) {
	args := m.Called(cvID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*analysis.AnalysisResult), args.Error(1)
}

// GetAnalysisResultsByJD mocks the GetAnalysisResultsByJD method
func (m *MockStorageService) GetAnalysisResultsByJD(jdID string) ([]*analysis.AnalysisResult, error) {
	args := m.Called(jdID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*analysis.AnalysisResult), args.Error(1)
}

var _ ports.StorageService = (*MockStorageService)(nil)
