package llm

import (
	// "cv-analyzer/internal/domain/analysis" // Removed as unused in these specific tests
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"os"
	"strings" // For Contains checks
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testOpenAIApiKey = "test-openai-api-key"

func TestNewOpenAIAdapter_WithAPIKey(t *testing.T) {
	adapter, err := NewOpenAIAdapter(testOpenAIApiKey)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)
	assert.Equal(t, testOpenAIApiKey, adapter.apiKey)
	assert.NotNil(t, adapter.client) // client is initialized in NewOpenAIAdapter
}

func TestNewOpenAIAdapter_WithEnvVar(t *testing.T) {
	originalApiKey := os.Getenv("OPENAI_API_KEY") // Preserve original value
	os.Setenv("OPENAI_API_KEY", testOpenAIApiKey)
	defer os.Setenv("OPENAI_API_KEY", originalApiKey) // Restore original value

	adapter, err := NewOpenAIAdapter("") // Pass empty string to force reading from env
	assert.NoError(t, err)
	assert.NotNil(t, adapter)
	assert.Equal(t, testOpenAIApiKey, adapter.apiKey)
	assert.NotNil(t, adapter.client)
}

func TestNewOpenAIAdapter_NoAPIKey_ShouldReturnError(t *testing.T) { // Renamed test
	// Ensure OPENAI_API_KEY is not set
	originalApiKey := os.Getenv("OPENAI_API_KEY")
	os.Unsetenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", originalApiKey) // Restore original value

	adapter, err := NewOpenAIAdapter("")
	assert.Error(t, err)
	assert.Nil(t, adapter)
	// NewOpenAIAdapter returns "OpenAI API key is not configured either directly or via OPENAI_API_KEY env var"
	assert.EqualError(t, err, "OpenAI API key is not configured either directly or via OPENAI_API_KEY env var")
}

func TestOpenAIAdapter_Analyze_Success_Mocked(t *testing.T) {
	adapter, err := NewOpenAIAdapter(testOpenAIApiKey)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	cvID := uuid.NewString()
	jdID := uuid.NewString()
	cvData := &cv.CV{ID: cvID, FilePath: "cv.pdf", TextContent: "Sample CV content"}
	jdData := &jd.JobDescription{ID: jdID, FilePath: "jd.pdf", TextContent: "Sample JD content"}

	result, err := adapter.Analyze(cvData, jdData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)

	// Assertions based on current OpenAIAdapter.Analyze mock implementation
	assert.Equal(t, cvID, result.CV.ID)
	assert.Equal(t, cvData.FilePath, result.CV.FilePath)
	assert.Equal(t, jdID, result.JobDescription.ID)
	assert.Equal(t, jdData.FilePath, result.JobDescription.FilePath)

	assert.Equal(t, 0.85, result.Score) // Current mock score is 0.85
	expectedSummaryPart := "Mock OpenAI Analysis for CV " + cvID
	assert.True(t, strings.Contains(result.Summary, expectedSummaryPart), "Summary should contain expected mock content")

	assert.Equal(t, []string{"example_missing_from_openai_adapter"}, result.MissingKeywords)
	assert.Equal(t, []string{"example_matching_from_openai_adapter"}, result.MatchingKeywords)
}

func TestOpenAIAdapter_Analyze_NilCV(t *testing.T) {
	adapter, _ := NewOpenAIAdapter(testOpenAIApiKey) // Assuming key is valid for constructor
	jdData := &jd.JobDescription{ID: uuid.NewString(), TextContent: "Sample JD content"}

	result, err := adapter.Analyze(nil, jdData)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "CV or JD data cannot be nil")
}

func TestOpenAIAdapter_Analyze_NilJD(t *testing.T) {
	adapter, _ := NewOpenAIAdapter(testOpenAIApiKey)
	cvData := &cv.CV{ID: uuid.NewString(), TextContent: "Sample CV content"}

	result, err := adapter.Analyze(cvData, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "CV or JD data cannot be nil")
}

func TestOpenAIAdapter_Analyze_EmptyTextContent_CV(t *testing.T) {
	adapter, _ := NewOpenAIAdapter(testOpenAIApiKey)
	cvData := &cv.CV{ID: uuid.NewString(), TextContent: ""} // Empty CV text
	jdData := &jd.JobDescription{ID: uuid.NewString(), TextContent: "Sample JD"}

	// OpenAIAdapter.Analyze has checks:
	// if cvData.TextContent == "" { return nil, errors.New("CV text content is empty") }
	// if jdData.TextContent == "" { return nil, errors.New("JD text content is empty") }
	result, err := adapter.Analyze(cvData, jdData)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "CV text content is empty")
}

func TestOpenAIAdapter_Analyze_EmptyTextContent_JD(t *testing.T) {
	adapter, _ := NewOpenAIAdapter(testOpenAIApiKey)
	cvData := &cv.CV{ID: uuid.NewString(), TextContent: "Sample CV"}
	jdData := &jd.JobDescription{ID: uuid.NewString(), TextContent: ""} // Empty JD text

	result, err := adapter.Analyze(cvData, jdData)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "JD text content is empty")
}
