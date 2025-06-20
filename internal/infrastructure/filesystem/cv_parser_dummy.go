package filesystem

import (
	"cv-analyzer/internal/domain/cv"
	"fmt"
	"github.com/google/uuid"
)

// DummyCVParser is a placeholder implementation for CVParser.
type DummyCVParser struct{}

// NewDummyCVParser creates a new DummyCVParser.
func NewDummyCVParser() *DummyCVParser {
	return &DummyCVParser{}
}

// ParseCV simulates parsing a CV file.
func (p *DummyCVParser) ParseCV(filePath string) (*cv.CV, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	// Simulate reading and parsing a CV
	// In a real implementation, this would involve opening the file,
	// reading its content, and extracting relevant information.
	fmt.Printf("Simulating parsing CV from: %s\n", filePath)
	// Provide dummy text content for the constructor
	dummyTextContent := fmt.Sprintf("This is dummy text content for CV file: %s", filePath)
	return cv.NewCV(uuid.NewString(), filePath, dummyTextContent), nil
}
