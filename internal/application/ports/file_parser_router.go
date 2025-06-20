package ports

// FileParserRouter defines the interface for a router that selects a file parser based on file type.
type FileParserRouter interface {
	GetParser(filePath string) (FileParser, error)
}
