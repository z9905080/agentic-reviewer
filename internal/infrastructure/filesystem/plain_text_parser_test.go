package filesystem

import (
	// "cv-analyzer/internal/domain/cv" // Not directly used if asserting on struct fields
	// "cv-analyzer/internal/domain/jd" // Not directly used
	"os" // For temp file creation
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainTextParser_ParseCV_Success(t *testing.T) {
	parser := NewPlainTextParser()
	// File content is not read from disk in this parser (this comment is outdated)

	// The PlainTextParser's ParseCV/ParseJD methods were refactored
	// to take filePath and then they read the file.
	// For unit testing, we should probably test the core logic (byte to string)
	// or use a temporary file.
	// Let's assume we want to test the logic that was:
	// Parse(fileContent []byte) (string, error)
	// And then ensure the CV/JD objects are built correctly.

	// Re-checking PlainTextParser's actual implementation:
	// It has `ParseCV(filePath string) (*cv.CV, error)`
	// which calls `os.ReadFile(filePath)`.
	// This makes it an integration test with the filesystem.
	// For a unit test, we'd ideally pass []byte.
	// Let's create a temporary file for this test.

	textContent := "This is a test CV content."
	tmpFile, err := os.CreateTemp("", "test_cv_*.txt")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	_, err = tmpFile.WriteString(textContent)
	assert.NoError(t, err)
	err = tmpFile.Close() // Close the file before parser tries to read it
	assert.NoError(t, err)


	parsedCV, err := parser.ParseCV(tmpFile.Name())

	assert.NoError(t, err)
	assert.NotNil(t, parsedCV)
	assert.Equal(t, tmpFile.Name(), parsedCV.FilePath)
	assert.Equal(t, textContent, parsedCV.TextContent)
	assert.NotEmpty(t, parsedCV.ID)
}

func TestPlainTextParser_ParseCV_FileNotExist(t *testing.T) {
	parser := NewPlainTextParser()
	filePath := "non_existent_cv.txt"

	parsedCV, err := parser.ParseCV(filePath)

	assert.Error(t, err) // Expect an error because file doesn't exist
	assert.Nil(t, parsedCV)
	assert.ErrorIs(t, err, os.ErrNotExist) // Check for specific error type if possible, or contains
}

func TestPlainTextParser_ParseCV_EmptyFile(t *testing.T) {
	parser := NewPlainTextParser()
	tmpFile, err := os.CreateTemp("", "empty_cv_*.txt")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	err = tmpFile.Close() // File is empty
	assert.NoError(t, err)

	// Depending on implementation, empty content might be an error or valid
	// Current PlainTextParser (if it uses the old logic) might error on empty string
	// The current ParseCV reads file, then creates CV. Empty content is fine.
	parsedCV, err := parser.ParseCV(tmpFile.Name())

	assert.NoError(t, err)
    assert.NotNil(t, parsedCV)
	assert.Equal(t, "", parsedCV.TextContent) // Expect empty text content
}


func TestPlainTextParser_ParseJD_Success(t *testing.T) {
	parser := NewPlainTextParser()
	textContent := "This is a test JD content."
	tmpFile, err := os.CreateTemp("", "test_jd_*.txt")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.WriteString(textContent)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)


	parsedJD, err := parser.ParseJD(tmpFile.Name())

	assert.NoError(t, err)
	assert.NotNil(t, parsedJD)
	assert.Equal(t, tmpFile.Name(), parsedJD.FilePath)
	assert.Equal(t, textContent, parsedJD.TextContent)
	assert.NotEmpty(t, parsedJD.ID)
}
