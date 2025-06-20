package usecases

import (
	// "cv-analyzer/internal/application/ports" // Not directly needed in test file, mocks use it
	"cv-analyzer/internal/application/usecases/mocks" // Our mocks
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"errors"
	"testing"
	// "time" // Not used if AnalyzedAt is not part of AnalysisResult domain object

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAnalyzeCVUseCase_Execute_Success(t *testing.T) {
	mockParserRouter := new(mocks.MockFileParserRouter)
	mockFileParser := new(mocks.MockFileParser) // This can be a single instance if GetParser returns it for both calls
	mockAnalysisService := new(mocks.MockAnalysisService)
	mockStorageService := new(mocks.MockStorageService)

	// Note: NewAnalyzeCVUseCase expects concrete type, not interface here based on current code.
	// If NewAnalyzeCVUseCase was changed to accept interfaces (which is good practice), this would be:
	// useCase := NewAnalyzeCVUseCase(mockParserRouter, mockAnalysisService, mockStorageService)
	// For now, assuming it's as previously defined:
	// func NewAnalyzeCVUseCase(parserRouter ports.FileParserRouter, analysisService ports.AnalysisService, storageService ports.StorageService)
	useCase := NewAnalyzeCVUseCase(mockParserRouter, mockAnalysisService, mockStorageService)

	cvFilePath := "path/to/cv.pdf"
	jdFilePath := "path/to/jd.txt"

	mockCV := &cv.CV{ID: uuid.NewString(), FilePath: cvFilePath, TextContent: "cv content"}
	mockJD := &jd.JobDescription{ID: uuid.NewString(), FilePath: jdFilePath, TextContent: "jd content"}

	// Adjusted mockAnalysisResult to match the domain.AnalysisResult struct
	mockAnalysisResult := &analysis.AnalysisResult{
		ID:               uuid.NewString(),
		CV:               mockCV,
		JobDescription:   mockJD,
		Score:            85.0,
		Summary:          "Good candidate",
		MatchingKeywords: []string{"Go", "Test"},
		MissingKeywords:  []string{"Docker"},
	}

	// Setup expectations
	mockParserRouter.On("GetParser", cvFilePath).Return(mockFileParser, nil)
	mockParserRouter.On("GetParser", jdFilePath).Return(mockFileParser, nil) // Assuming same parser for .txt and .pdf for this mock setup

	mockFileParser.On("ParseCV", cvFilePath).Return(mockCV, nil)
	mockFileParser.On("ParseJD", jdFilePath).Return(mockJD, nil)

	mockAnalysisService.On("Analyze", mockCV, mockJD).Return(mockAnalysisResult, nil)

	// Mocking storage calls made by the use case
	mockStorageService.On("SaveCV", mockCV).Return(nil) // Expect SaveCV to be called
	mockStorageService.On("SaveJD", mockJD).Return(nil) // Expect SaveJD to be called
	mockStorageService.On("SaveAnalysisResult", mockAnalysisResult).Return(nil)


	result, err := useCase.Execute(cvFilePath, jdFilePath)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockAnalysisResult.ID, result.ID)
	assert.Equal(t, mockAnalysisResult.Score, result.Score)
	assert.Equal(t, mockCV.ID, result.CV.ID) // Verify nested data
	assert.Equal(t, mockJD.ID, result.JobDescription.ID)

	mockParserRouter.AssertExpectations(t)
	mockFileParser.AssertExpectations(t)
	mockAnalysisService.AssertExpectations(t)
	mockStorageService.AssertExpectations(t)
}

func TestAnalyzeCVUseCase_Execute_ParserRouterError_CV(t *testing.T) {
	mockParserRouter := new(mocks.MockFileParserRouter)
	mockAnalysisService := new(mocks.MockAnalysisService)
	mockStorageService := new(mocks.MockStorageService)

	useCase := NewAnalyzeCVUseCase(mockParserRouter, mockAnalysisService, mockStorageService)
	cvFilePath := "path/to/cv.unsupported"
	jdFilePath := "path/to/jd.txt" // This won't be used if CV parsing selection fails
	expectedError := errors.New("unsupported file type for CV")

	mockParserRouter.On("GetParser", cvFilePath).Return(nil, expectedError)
	// GetParser for JD should not be called if getting CV parser fails.

	result, err := useCase.Execute(cvFilePath, jdFilePath)

	assert.Error(t, err)
	assert.Nil(t, result)
	// Error message comes from AnalyzeCVUseCase: fmt.Errorf("failed to get parser for CV file %s: %w", cvFilePath, err)
	assert.Contains(t, err.Error(), "failed to get parser for CV file path/to/cv.unsupported")
	assert.Contains(t, err.Error(), "unsupported file type for CV")


	mockParserRouter.AssertExpectations(t)
	mockAnalysisService.AssertNotCalled(t, "Analyze", mock.Anything, mock.Anything)
	mockStorageService.AssertNotCalled(t, "SaveAnalysisResult", mock.Anything)
	mockStorageService.AssertNotCalled(t, "SaveCV", mock.Anything)
	mockStorageService.AssertNotCalled(t, "SaveJD", mock.Anything)
}

func TestAnalyzeCVUseCase_Execute_CVParseError(t *testing.T) {
	mockParserRouter := new(mocks.MockFileParserRouter)
	mockFileParser := new(mocks.MockFileParser)
	mockAnalysisService := new(mocks.MockAnalysisService)
	mockStorageService := new(mocks.MockStorageService)

	useCase := NewAnalyzeCVUseCase(mockParserRouter, mockAnalysisService, mockStorageService)
	cvFilePath := "path/to/cv.pdf"
	jdFilePath := "path/to/jd.txt"
	expectedError := errors.New("cv parse failed")

	mockParserRouter.On("GetParser", cvFilePath).Return(mockFileParser, nil)
	mockFileParser.On("ParseCV", cvFilePath).Return(nil, expectedError)
	// GetParser for JD and subsequent calls should not happen if ParseCV fails.

	result, err := useCase.Execute(cvFilePath, jdFilePath)

	assert.Error(t, err)
	assert.Nil(t, result)
	// Error message from AnalyzeCVUseCase: fmt.Errorf("failed to parse CV file %s: %w", cvFilePath, err)
	assert.Contains(t, err.Error(), "failed to parse CV file path/to/cv.pdf")
	assert.Contains(t, err.Error(), "cv parse failed")


	mockParserRouter.AssertExpectations(t) // Checks GetParser(cvFilePath)
	mockFileParser.AssertExpectations(t)   // Checks ParseCV(cvFilePath)
	mockParserRouter.AssertNotCalled(t, "GetParser", jdFilePath) // Should not be called
	mockAnalysisService.AssertNotCalled(t, "Analyze", mock.Anything, mock.Anything)
}


func TestAnalyzeCVUseCase_Execute_AnalysisServiceError(t *testing.T) {
    mockParserRouter := new(mocks.MockFileParserRouter)
    mockFileParser := new(mocks.MockFileParser)
    mockAnalysisService := new(mocks.MockAnalysisService)
    mockStorageService := new(mocks.MockStorageService)

    useCase := NewAnalyzeCVUseCase(mockParserRouter, mockAnalysisService, mockStorageService)

    cvFilePath := "path/to/cv.pdf"
    jdFilePath := "path/to/jd.txt"

    mockCV := &cv.CV{ID: uuid.NewString(), FilePath: cvFilePath, TextContent: "cv content"}
    mockJD := &jd.JobDescription{ID: uuid.NewString(), FilePath: jdFilePath, TextContent: "jd content"}
    analysisError := errors.New("llm analysis failed")

    mockParserRouter.On("GetParser", cvFilePath).Return(mockFileParser, nil)
    mockParserRouter.On("GetParser", jdFilePath).Return(mockFileParser, nil)
    mockFileParser.On("ParseCV", cvFilePath).Return(mockCV, nil)
    mockFileParser.On("ParseJD", jdFilePath).Return(mockJD, nil)
    mockAnalysisService.On("Analyze", mockCV, mockJD).Return(nil, analysisError)
	// SaveCV and SaveJD should NOT be called if Analyze fails.
    // mockStorageService.On("SaveCV", mockCV).Return(nil)
    // mockStorageService.On("SaveJD", mockJD).Return(nil)


    result, err := useCase.Execute(cvFilePath, jdFilePath)

    assert.Error(t, err)
    assert.Nil(t, result)
	// Error from AnalyzeCVUseCase: fmt.Errorf("failed to analyze CV and JD: %w", err)
    assert.Contains(t, err.Error(), "failed to analyze CV and JD")
    assert.Contains(t, err.Error(), "llm analysis failed")

    mockParserRouter.AssertExpectations(t)
    mockFileParser.AssertExpectations(t)
    mockAnalysisService.AssertExpectations(t)
    mockStorageService.AssertNotCalled(t, "SaveCV", mockCV) // Should not be called
    mockStorageService.AssertNotCalled(t, "SaveJD", mockJD) // Should not be called
    mockStorageService.AssertNotCalled(t, "SaveAnalysisResult", mock.Anything)
}

func TestAnalyzeCVUseCase_Execute_StorageServiceError_SaveAnalysisResult(t *testing.T) {
    mockParserRouter := new(mocks.MockFileParserRouter)
    mockFileParser := new(mocks.MockFileParser)
    mockAnalysisService := new(mocks.MockAnalysisService)
    mockStorageService := new(mocks.MockStorageService)

    useCase := NewAnalyzeCVUseCase(mockParserRouter, mockAnalysisService, mockStorageService)

    cvFilePath := "path/to/cv.pdf"
    jdFilePath := "path/to/jd.txt"

    mockCV := &cv.CV{ID: uuid.NewString(), FilePath: cvFilePath, TextContent: "cv content"}
    mockJD := &jd.JobDescription{ID: uuid.NewString(), FilePath: jdFilePath, TextContent: "jd content"}
    mockAnalysisResult := &analysis.AnalysisResult{
		ID: uuid.NewString(), CV: mockCV, JobDescription: mockJD, Score: 75.0, Summary: "Summary",
	}
    storageError := errors.New("db save analysis failed")


    mockParserRouter.On("GetParser", cvFilePath).Return(mockFileParser, nil)
    mockParserRouter.On("GetParser", jdFilePath).Return(mockFileParser, nil)
    mockFileParser.On("ParseCV", cvFilePath).Return(mockCV, nil)
    mockFileParser.On("ParseJD", jdFilePath).Return(mockJD, nil)
    mockAnalysisService.On("Analyze", mockCV, mockJD).Return(mockAnalysisResult, nil)
    mockStorageService.On("SaveCV", mockCV).Return(nil) // Successful
    mockStorageService.On("SaveJD", mockJD).Return(nil) // Successful
    mockStorageService.On("SaveAnalysisResult", mockAnalysisResult).Return(storageError) // This one fails

    result, err := useCase.Execute(cvFilePath, jdFilePath)

    assert.Error(t, err)
    assert.Nil(t, result) // Because SaveAnalysisResult error path returns (nil, err)
	// Error from AnalyzeCVUseCase: fmt.Errorf("failed to save analysis result: %w", err)
    assert.Contains(t, err.Error(), "failed to save analysis result")
    assert.Contains(t, err.Error(), "db save analysis failed")


    mockParserRouter.AssertExpectations(t)
    mockFileParser.AssertExpectations(t)
    mockAnalysisService.AssertExpectations(t)
    mockStorageService.AssertExpectations(t)
}

// Test for when SaveCV fails (should log warning, not return error)
func TestAnalyzeCVUseCase_Execute_StorageServiceError_SaveCV_Warning(t *testing.T) {
	mockParserRouter := new(mocks.MockFileParserRouter)
	mockFileParser := new(mocks.MockFileParser)
	mockAnalysisService := new(mocks.MockAnalysisService)
	mockStorageService := new(mocks.MockStorageService)

	useCase := NewAnalyzeCVUseCase(mockParserRouter, mockAnalysisService, mockStorageService)

	cvFilePath := "path/to/cv.pdf"
	jdFilePath := "path/to/jd.txt"

	mockCV := &cv.CV{ID: uuid.NewString(), FilePath: cvFilePath, TextContent: "cv content"}
	mockJD := &jd.JobDescription{ID: uuid.NewString(), FilePath: jdFilePath, TextContent: "jd content"}
	mockAnalysisResult := &analysis.AnalysisResult{ID: uuid.NewString(), CV: mockCV, JobDescription: mockJD, Score: 90.0}
	saveCVErr := errors.New("failed to save CV")

	mockParserRouter.On("GetParser", cvFilePath).Return(mockFileParser, nil)
	mockParserRouter.On("GetParser", jdFilePath).Return(mockFileParser, nil)
	mockFileParser.On("ParseCV", cvFilePath).Return(mockCV, nil)
	mockFileParser.On("ParseJD", jdFilePath).Return(mockJD, nil)
	mockAnalysisService.On("Analyze", mockCV, mockJD).Return(mockAnalysisResult, nil)

	mockStorageService.On("SaveCV", mockCV).Return(saveCVErr) // This fails
	mockStorageService.On("SaveJD", mockJD).Return(nil)      // This succeeds
	mockStorageService.On("SaveAnalysisResult", mockAnalysisResult).Return(nil) // This succeeds

	result, err := useCase.Execute(cvFilePath, jdFilePath)

	assert.NoError(t, err) // No error returned to caller for SaveCV failure
	assert.NotNil(t, result)
	assert.Equal(t, mockAnalysisResult.ID, result.ID)

	// Verify all mocks were called as expected
	mockParserRouter.AssertExpectations(t)
	mockFileParser.AssertExpectations(t)
	mockAnalysisService.AssertExpectations(t)
	mockStorageService.AssertExpectations(t)
	// Further: check log output for warning (requires log mocking or capturing stdout, more advanced)
}
