package ports

import "cv-analyzer/internal/domain/jd"

// JDParser defines the interface for parsing Job Description files.
type JDParser interface {
	ParseJD(filePath string) (*jd.JobDescription, error)
}
