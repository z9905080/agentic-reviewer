package mocks

import (
	"cv-analyzer/internal/application/ports"
	"github.com/stretchr/testify/mock"
)

type MockFileParserRouter struct {
	mock.Mock
}

func (m *MockFileParserRouter) GetParser(filePath string) (ports.FileParser, error) {
	args := m.Called(filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(ports.FileParser), args.Error(1)
}
var _ ports.FileParserRouter = (*MockFileParserRouter)(nil)
