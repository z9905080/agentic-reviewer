package filesystem

import (
	"cv-analyzer/internal/application/ports" // For FileParser interface
	"cv-analyzer/internal/domain/cv"         // Added for MinimalMockFileParser method signatures
	"cv-analyzer/internal/domain/jd"         // Added for MinimalMockFileParser method signatures
	"strings"                                // Added for error checking in test
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MinimalMockFileParser for ParserRouter tests (if we wanted to mock internal parsers)
// However, the current test setup uses real instances of PlainTextParser and PdfParser.
// This mock is not strictly needed for the current version of TestParserRouter_GetParser_Success
// but might be useful for other scenarios or if NewParserRouter was mocked.
type MinimalMockFileParser struct {
	mock.Mock
	Type string // To identify which mock parser it is
}
func (m *MinimalMockFileParser) ParseCV(filePath string) (*cv.CV, error) {
	args := m.Called(filePath)
	if args.Get(0) == nil && args.Error(1) == nil { // Handle case where both are nil (e.g. if Get(0) isn't set for a success with nil CV)
        return nil, nil
    }
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
	return args.Get(0).(*cv.CV), args.Error(1)
}
func (m *MinimalMockFileParser) ParseJD(filePath string) (*jd.JobDescription, error) {
	args := m.Called(filePath)
    if args.Get(0) == nil && args.Error(1) == nil {
        return nil, nil
    }
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
	return args.Get(0).(*jd.JobDescription), args.Error(1)
}
var _ ports.FileParser = (*MinimalMockFileParser)(nil)


func TestParserRouter_GetParser_Success(t *testing.T) {
	// NewParserRouter was refactored to take its parsers for DI.
	realTxtParser := NewPlainTextParser()
	realPdfParser := NewPdfParser()

	router, err := NewParserRouter(realTxtParser, realPdfParser)
	assert.NoError(t, err)
	assert.NotNil(t, router)

	testCases := []struct {
		name         string
		filePath     string
		expectedType ports.FileParser
		expectError  bool
	}{
		{"Text File", "document.txt", realTxtParser, false},
		{"PDF File", "document.pdf", realPdfParser, false},
		{"Uppercase PDF", "DOCUMENT.PDF", realPdfParser, false},
		{"Text File with long path", "/long/path/to/my/cv.txt", realTxtParser, false},
		{"Unsupported Extension", "document.docx", nil, true},
		{"No Extension", "document", nil, true},
		{"Only Dot", ".", nil, true},
		{"Empty Filename", "", nil, true}, // Added test case for empty filename
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser, err := router.GetParser(tc.filePath)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, parser)
				// Check for specific error messages based on the type of error expected
				if tc.filePath == "" {
					// filepath.Ext("") is "", so it hits the default "unsupported file type: "
					assert.Contains(t, err.Error(), "unsupported file type: ")
				} else if strings.HasSuffix(tc.filePath, ".docx") || (tc.filePath == "document" && !strings.Contains(tc.filePath, ".")) || tc.filePath == "." {
					assert.Contains(t, err.Error(), "unsupported file type")
				}
				// Other specific error checks could be added here
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, parser)
				assert.Same(t, tc.expectedType, parser, "Parser instance mismatch for %s", tc.filePath)
			}
		})
	}
}

func TestParserRouter_NewParserRouter_NilDependencies(t *testing.T) {
    // Test NewParserRouter's error handling for nil dependencies.
    _, err := NewParserRouter(nil, NewPdfParser())
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "PlainTextParser cannot be nil")

    _, err = NewParserRouter(NewPlainTextParser(), nil)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "PdfParser cannot be nil")
}

// Test GetParser with empty filepath specifically if not covered by table test's error message.
func TestParserRouter_GetParser_EmptyFilePath(t *testing.T) {
	realTxtParser := NewPlainTextParser()
	realPdfParser := NewPdfParser()
	router, _ := NewParserRouter(realTxtParser, realPdfParser)

	parser, err := router.GetParser("")
	assert.Error(t, err)
	assert.Nil(t, parser)
	// Add a specific check for empty filepath error message if GetParser implements it.
	// For now, the default "unsupported file type: " might be returned by the switch's default.
	// Let's assume filepath.Ext("") is "", so it hits the default.
	assert.Contains(t, err.Error(), "unsupported file type: ")
	// If we want a more specific error for empty string before extension check,
	// GetParser would need an initial check like:
	// if strings.TrimSpace(filePath) == "" { return nil, fmt.Errorf("file path cannot be empty in GetParser") }
}
