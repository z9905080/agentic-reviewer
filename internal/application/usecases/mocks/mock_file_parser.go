package mocks

import (
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"cv-analyzer/internal/application/ports" // For FileParser interface
	"github.com/stretchr/testify/mock"
)

type MockFileParser struct {
	mock.Mock
}

func (m *MockFileParser) ParseCV(filePath string) (*cv.CV, error) {
	args := m.Called(filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cv.CV), args.Error(1)
}

func (m *MockFileParser) ParseJD(filePath string) (*jd.JobDescription, error) {
	args := m.Called(filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jd.JobDescription), args.Error(1)
}

var _ ports.FileParser = (*MockFileParser)(nil)
