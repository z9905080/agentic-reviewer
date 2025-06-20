package filesystem

import (
	"cv-analyzer/internal/application/ports"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"errors"
	"fmt" // Added missing import
	"github.com/google/uuid"
	"os"      // For reading files
	"strings" // For basic text manipulation
)

// PlainTextParser implements CVParser and JDParser for plain text files.
type PlainTextParser struct{}

// NewPlainTextParser creates a new PlainTextParser instance.
func NewPlainTextParser() *PlainTextParser {
	return &PlainTextParser{}
}

// ParseCV reads a CV file from filePath, extracts its text content,
// and returns a new CV domain object.
func (p *PlainTextParser) ParseCV(filePath string) (*cv.CV, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, errors.New("file path cannot be empty")
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CV file %s: %w", filePath, err)
	}

	textContent := string(fileContent)
	// Optional: Add validation for empty content if necessary
	// if strings.TrimSpace(textContent) == "" {
	// 	return nil, errors.New("cv file content is empty or only whitespace")
	// }

	newID := uuid.NewString()
	cvInstance := cv.NewCV(newID, filePath, textContent)

	return cvInstance, nil
}

// ParseJD reads a Job Description file from filePath, extracts its text content,
// and returns a new JobDescription domain object.
func (p *PlainTextParser) ParseJD(filePath string) (*jd.JobDescription, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, errors.New("file path cannot be empty")
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JD file %s: %w", filePath, err)
	}

	textContent := string(fileContent)
	// Optional: Add validation for empty content if necessary
	// if strings.TrimSpace(textContent) == "" {
	// 	return nil, errors.New("jd file content is empty or only whitespace")
	// }

	newID := uuid.NewString()
	jdInstance := jd.NewJobDescription(newID, filePath, textContent)

	return jdInstance, nil
}

// Compile-time checks to ensure PlainTextParser implements the interfaces.
var _ ports.CVParser = (*PlainTextParser)(nil)
var _ ports.JDParser = (*PlainTextParser)(nil)
