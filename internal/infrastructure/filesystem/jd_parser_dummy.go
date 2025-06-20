package filesystem

import (
	"cv-analyzer/internal/domain/jd"
	"fmt"
	"github.com/google/uuid"
)

// DummyJDParser is a placeholder implementation for JDParser.
type DummyJDParser struct{}

// NewDummyJDParser creates a new DummyJDParser.
func NewDummyJDParser() *DummyJDParser {
	return &DummyJDParser{}
}

// ParseJD simulates parsing a Job Description file.
func (p *DummyJDParser) ParseJD(filePath string) (*jd.JobDescription, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	// Simulate reading and parsing a JD
	fmt.Printf("Simulating parsing Job Description from: %s\n", filePath)
	// Provide dummy text content for the constructor
	dummyTextContent := fmt.Sprintf("This is dummy text content for JD file: %s", filePath)
	return jd.NewJobDescription(uuid.NewString(), filePath, dummyTextContent), nil
}
