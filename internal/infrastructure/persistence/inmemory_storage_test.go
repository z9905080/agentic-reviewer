package persistence

import (
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv" // For nested CV in AnalysisResult
	"cv-analyzer/internal/domain/jd" // For nested JD in AnalysisResult
	"testing"
	// "time" // Not strictly needed if we don't assert on it directly for these tests

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStorage_SaveAndGetAnalysisResult_Success(t *testing.T) {
	storage, err := NewInMemoryStorage()
	assert.NoError(t, err)

	resultID := uuid.NewString()
	mockCV := &cv.CV{ID: uuid.NewString(), FilePath: "cv.txt", TextContent: "cv text"}
	mockJD := &jd.JobDescription{ID: uuid.NewString(), FilePath: "jd.txt", TextContent: "jd text"}

	analysisResult := &analysis.AnalysisResult{
		ID:               resultID,
		CV:               mockCV,
		JobDescription:   mockJD,
		Score:            90.5,
		Summary:          "Excellent candidate",
		MatchingKeywords: []string{"Go", "Microservices"},
		MissingKeywords:  []string{"Java"},
		// AnalyzedAt:       time.Now(), // Not asserting this directly
	}

	err = storage.SaveAnalysisResult(analysisResult) // Using the correct method name
	assert.NoError(t, err)

	retrievedResult, err := storage.GetAnalysisResult(resultID) // Using the correct method name
	assert.NoError(t, err)
	assert.NotNil(t, retrievedResult)
	assert.Equal(t, analysisResult.ID, retrievedResult.ID)
	assert.Equal(t, analysisResult.Summary, retrievedResult.Summary)
	assert.Equal(t, analysisResult.Score, retrievedResult.Score)
	assert.Equal(t, analysisResult.CV.ID, retrievedResult.CV.ID) // Check nested data
}

func TestInMemoryStorage_GetAnalysisResult_NotFound(t *testing.T) {
	storage, err := NewInMemoryStorage()
	assert.NoError(t, err)

	nonExistentID := uuid.NewString()
	retrievedResult, err := storage.GetAnalysisResult(nonExistentID) // Using the correct method name

	assert.Error(t, err)
	assert.Nil(t, retrievedResult)
	assert.Contains(t, err.Error(), "not found")
}

func TestInMemoryStorage_SaveAnalysisResult_AlreadyExists(t *testing.T) {
	storage, err := NewInMemoryStorage()
	assert.NoError(t, err)

	resultID := uuid.NewString()
	analysisResult1 := &analysis.AnalysisResult{ID: resultID, Summary: "First save"}

	err = storage.SaveAnalysisResult(analysisResult1) // Using the correct method name
	assert.NoError(t, err)

	analysisResult2 := &analysis.AnalysisResult{ID: resultID, Summary: "Second save attempt"}
	err = storage.SaveAnalysisResult(analysisResult2) // Using the correct method name

	assert.Error(t, err) // Now expecting an error due to the updated SaveAnalysisResult logic
	assert.Contains(t, err.Error(), "already exists")

	// Verify that the first one is still there and not overwritten
	retrievedResult, errGet := storage.GetAnalysisResult(resultID)
	assert.NoError(t, errGet)
	assert.NotNil(t, retrievedResult)
	assert.Equal(t, "First save", retrievedResult.Summary)
}


func TestInMemoryStorage_SaveAnalysisResult_NilResult(t *testing.T) {
    storage, err := NewInMemoryStorage()
    assert.NoError(t, err)

    // The current SaveAnalysisResult method checks for result == nil and result.ID == ""
	// if result == nil { return fmt.Errorf("analysis result cannot be nil") }
	// if result.ID == "" { return fmt.Errorf("analysis result ID cannot be empty") }
    // So, saving a nil result should return an error, not panic.
    err = storage.SaveAnalysisResult(nil)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "analysis result cannot be nil")
}

func TestInMemoryStorage_SaveAnalysisResult_EmptyID(t *testing.T) {
	storage, err := NewInMemoryStorage()
	assert.NoError(t, err)

	analysisResult := &analysis.AnalysisResult{ID: "", Summary: "No ID"}
	err = storage.SaveAnalysisResult(analysisResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "analysis result ID cannot be empty")
}

// TODO: Add tests for SaveCV, GetCV, SaveJD, GetJD, GetAnalysisResultsByCV, GetAnalysisResultsByJD
