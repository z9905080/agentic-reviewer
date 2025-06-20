package ports

import "cv-analyzer/internal/domain/cv"

// CVParser defines the interface for parsing CV files.
type CVParser interface {
	ParseCV(filePath string) (*cv.CV, error)
}
