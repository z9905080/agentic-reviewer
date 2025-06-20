package cv

// CV represents a curriculum vitae.
type CV struct {
	ID          string
	FilePath    string
	TextContent string // Holds the parsed text content of the CV
	// TODO: Add more fields like ParsedData (structured), etc.
}

// NewCV creates a new CV instance.
func NewCV(id string, filePath string, textContent string) *CV {
	return &CV{
		ID:       id,
		FilePath: filePath,
		TextContent: textContent,
	}
}
