package jd

// JobDescription represents a job description.
type JobDescription struct {
	ID          string
	FilePath    string
	TextContent string // Holds the parsed text content of the JD
	// TODO: Add more fields like ParsedData (structured), etc.
}

// NewJobDescription creates a new JobDescription instance.
func NewJobDescription(id string, filePath string, textContent string) *JobDescription {
	return &JobDescription{
		ID:          id,
		FilePath:    filePath,
		TextContent: textContent,
	}
}
