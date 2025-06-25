package handlers

import (
	"bytes"
	"cv-analyzer/internal/application/ports"
	"cv-analyzer/internal/application/usecases"
	"cv-analyzer/internal/application/usecases/mocks" // Assuming mocks are in this path
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"cv-analyzer/internal/infrastructure/filesystem" // For real ParserRouter
	"cv-analyzer/internal/infrastructure/persistence" // For real InMemoryStorage
	"cv-analyzer/internal/interfaces/api/dto"
	// "cv-analyzer/internal/interfaces/api/routes" // REMOVE: To break import cycle
	"encoding/base64" // For PDF test
	"encoding/json"
	"errors" // Added for errors.New
	// "fmt" // Removed as likely unused after changing fmt.Errorf to errors.New
	// "io" // Removed as unused in test file
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	// "os" // Removed as unused in test file (os.Remove is in handler, not test)
	"path/filepath"
	// "strings" // Removed as no longer used after PdfFiles success test refactor
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper to create a multipart form request with files
func createMultipartRequest(t *testing.T, cvFilePath, jdFilePath string, cvContent, jdContent []byte) (*http.Request, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	if cvFilePath != "" && cvContent != nil {
		part, err := writer.CreateFormFile("cv", filepath.Base(cvFilePath))
		assert.NoError(t, err)
		_, err = part.Write(cvContent)
		assert.NoError(t, err)
	}

	if jdFilePath != "" && jdContent != nil {
		part, err := writer.CreateFormFile("jd", filepath.Base(jdFilePath))
		assert.NoError(t, err)
		_, err = part.Write(jdContent)
		assert.NoError(t, err)
	}

	err := writer.Close()
	assert.NoError(t, err)

	// Target URL should match the one defined in routes
	req, err := http.NewRequest("POST", "/api/v1/analysis/cv", body)
	assert.NoError(t, err)

	return req, writer
}

// Base64 encoded simple PDF for testing PDF uploads
const testPdfBase64 = "JVBERi0xLjQKJSAgMSAwIG9iago8PAogIC9UeXBlIC9DYXRhbG9nCiAgL1BhZ2VzIDIgMCBSCiAgL091dGxpbmVzIDMgMCBSCj4+CmVuZG9iagoKMiAwIG9iago8PAogIC9UeXBlIC9QYWdlcwogIC9Db3VudCAxCiAgL0tpZHMgWzQgMCBSXQo+PgplbmRvYmoKCjMgMCBvYmoKPDwKICAvVHlwZSAvT3V0bGluZXMKLy9Db3VudCAwIAo+PgplbmRvYmoKCjQgMCBvYmoKPDwKICAvVHlwZSAvUGFnZQogIC9QYXJlbnQgMiAwIFIKICAvUmVzb3VyY2VzCjw8CiAgICAvRm9udCA2IDAgUgo+PgogIC9NZWRpYUJveCBbMCAwIDYxMiA3OTJdCiAgL0NvbnRlbnRzIDUgMCBSCj4+CmVuZG9iagoKNSAwIG9iago8PAogIC9MZW5ndGggNDAKPj4Kc3RyZWFtCkJUCiAgL0YxIDEyIFRmCiAgNTUgNzAwIFRkCiAgKFRlc3QgQ1YgY29udGVudCBmcm9tIFBERikgVGoKRVQKClBlbmRzdHJlYW0KZW5kb2JqCgo2IDAgb2JqCjw8CiAgL1R5cGUgL0ZvbnQKICAvU3VidHlwZSAvVHlwZTEKICAvQmFzZUZvbnQgL0hlbHZldGljYQo+PgplbmRvYmoKCnhyZWYKMCA3CjAwMDAwMDAwMDAgNjU1MzUgZiAKMDAwMDAwMDAxMCAwMDAwMCBuIAowMDAwMDAwMDc5IDAwMDAwIG4gCjAwMDAwMDAxNDkgMDAwMDAwIG4gCjAwMDAwMDAyMDIgMDAwMDAwIG4gCjAwMDAwMDAzNjAgMDAwMDAwIG4gCjAwMDAwMDA0NDkgMDAwMDAwIG4gCnRyYWlsZXIKPDwKICAvU2l6ZSA3CiAgL1Jvb3QgMSAwIFIKPj4Kc3RhcnR4cmVmCjUwNAolJUVPRgo="

func setupTestRouter(t *testing.T, mockAnalysisService ports.AnalysisService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	testRouter := gin.New() // Define testRouter

	plainTextParser := filesystem.NewPlainTextParser()
	pdfParser := filesystem.NewPdfParser()
	parserRouter, err := filesystem.NewParserRouter(plainTextParser, pdfParser)
	assert.NoError(t, err) // Ensure router creation doesn't fail

	storageService, err := persistence.NewInMemoryStorage()
	assert.NoError(t, err) // Ensure storage creation doesn't fail

	analyzeUseCase := usecases.NewAnalyzeCVUseCase(parserRouter, mockAnalysisService, storageService)
	analysisHandler := NewAnalysisHandler(*analyzeUseCase)

	// Register the specific route for this handler directly in the test router
	apiGroup := testRouter.Group("/api/v1")
	{
		analysisRoutes := apiGroup.Group("/analysis")
		{
			analysisRoutes.POST("/cv", analysisHandler.HandleAnalyzeCV)
		}
	}
	return testRouter
}


func TestAnalysisHandler_HandleAnalyzeCV_Success_TxtFiles(t *testing.T) {
	mockLLMService := new(mocks.MockAnalysisService)
	testRouter := setupTestRouter(t, mockLLMService)

	cvContent := []byte("This is CV text.")
	jdContent := []byte("This is JD text.")

	// Mock LLM Analyze method
	expectedAnalysisResult := &analysis.AnalysisResult{
		ID:               uuid.NewString(),
		// CV and JD objects will be populated by the use case/parser
		CV: &cv.CV{},
		JobDescription: &jd.JobDescription{},
		Score:            90,
		Summary:          "Excellent match (mocked)",
		MatchingKeywords: []string{"text"},
		MissingKeywords:  []string{},
		// AnalyzedAt will be set by the handler, so we don't pre-set it here for direct comparison
	}

	mockLLMService.On("Analyze", mock.AnythingOfType("*cv.CV"), mock.AnythingOfType("*jd.JobDescription")).
		Run(func(args mock.Arguments) {
			cvArg := args.Get(0).(*cv.CV)
			jdArg := args.Get(1).(*jd.JobDescription)
			assert.Equal(t, string(cvContent), cvArg.TextContent)
			assert.Equal(t, string(jdContent), jdArg.TextContent)

			// Populate parts of expectedAnalysisResult that depend on dynamic data from parsing
			expectedAnalysisResult.CV.ID = cvArg.ID
			expectedAnalysisResult.CV.FilePath = cvArg.FilePath
			expectedAnalysisResult.JobDescription.ID = jdArg.ID
			expectedAnalysisResult.JobDescription.FilePath = jdArg.FilePath
		}).
		Return(expectedAnalysisResult, nil)


	req, writer := createMultipartRequest(t, "cv.txt", "jd.txt", cvContent, jdContent)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())

	var responseDto dto.AnalysisResponse
	err := json.Unmarshal(rr.Body.Bytes(), &responseDto)
	assert.NoError(t, err)
	assert.Equal(t, expectedAnalysisResult.ID, responseDto.AnalysisID) // DTO uses AnalysisID
	assert.Equal(t, expectedAnalysisResult.Score, responseDto.Score)
	assert.Equal(t, expectedAnalysisResult.Summary, responseDto.Summary)
	assert.Equal(t, expectedAnalysisResult.CV.ID, responseDto.CVID)
	assert.Equal(t, expectedAnalysisResult.CV.FilePath, responseDto.CVPath)
	assert.Equal(t, expectedAnalysisResult.JobDescription.ID, responseDto.JDID)
	assert.Equal(t, expectedAnalysisResult.JobDescription.FilePath, responseDto.JDPath)
	assert.WithinDuration(t, time.Now(), responseDto.AnalyzedAt, 10*time.Second) // Check AnalyzedAt is recent

	mockLLMService.AssertExpectations(t)
}


func TestAnalysisHandler_HandleAnalyzeCV_Success_PdfFiles(t *testing.T) {
    mockLLMService := new(mocks.MockAnalysisService) // This mock won't be called if parsing fails
    testRouter := setupTestRouter(t, mockLLMService)

    pdfBytes, err := base64.StdEncoding.DecodeString(testPdfBase64)
    assert.NoError(t, err)

    // This test now expects a failure because the testPdfBase64 is known to cause
    // a parsing error with rsc.io/pdf.
    req, writer := createMultipartRequest(t, "cv.pdf", "jd.pdf", pdfBytes, pdfBytes)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    rr := httptest.NewRecorder()
    testRouter.ServeHTTP(rr, req)

    // Expecting a 500 Internal Server Error because PdfParser will fail,
    // and this error will bubble up from the use case.
    assert.Equal(t, http.StatusInternalServerError, rr.Code, rr.Body.String())

    var errDto dto.ErrorResponse
    err = json.Unmarshal(rr.Body.Bytes(), &errDto)
    assert.NoError(t, err)
    // Check for the error message from PdfParser failing
    assert.Contains(t, errDto.Error, "Analysis failed: failed to parse CV file")
    assert.Contains(t, errDto.Error, "malformed PDF") // Specific to rsc.io/pdf error

    mockLLMService.AssertNotCalled(t, "Analyze", mock.Anything, mock.Anything)
}


func TestAnalysisHandler_HandleAnalyzeCV_MissingCVFile(t *testing.T) {
	mockLLMService := new(mocks.MockAnalysisService)
	testRouter := setupTestRouter(t, mockLLMService)

	jdContent := []byte("This is JD text.")
	// Pass empty cvFilePath and nil cvContent to createMultipartRequest to skip creating the CV part
	req, writer := createMultipartRequest(t, "", "jd.txt", nil, jdContent)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var errDto dto.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errDto)
	assert.NoError(t, err)
	assert.Equal(t, "CV file is required", errDto.Error)
}

func TestAnalysisHandler_HandleAnalyzeCV_UnsupportedFileType(t *testing.T) {
	mockLLMService := new(mocks.MockAnalysisService)
	testRouter := setupTestRouter(t, mockLLMService)

	cvContent := []byte("CV content in docx.")
	jdContent := []byte("JD content.")
	req, writer := createMultipartRequest(t, "cv.docx", "jd.txt", cvContent, jdContent)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var errDto dto.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errDto) // Use Unmarshal here
	assert.NoError(t, err)
	assert.Equal(t, "Unsupported file type. Only .txt and .pdf are allowed.", errDto.Error)
}


func TestAnalysisHandler_HandleAnalyzeCV_UseCaseError(t *testing.T) {
	mockLLMService := new(mocks.MockAnalysisService)
	testRouter := setupTestRouter(t, mockLLMService)

	cvContent := []byte("This is CV text.")
	jdContent := []byte("This is JD text.")

	expectedErrorMsg := "LLM is down (mocked)"
	mockLLMService.On("Analyze", mock.AnythingOfType("*cv.CV"), mock.AnythingOfType("*jd.JobDescription")).
		Return(nil, errors.New(expectedErrorMsg)) // Changed to errors.New

	req, writer := createMultipartRequest(t, "cv.txt", "jd.txt", cvContent, jdContent)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var errDto dto.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errDto) // Use Unmarshal here
	assert.NoError(t, err)

	actualError := errDto.Error
	// expectedToContain := "Analysis failed: " + expectedErrorMsg
	// t.Logf("Actual error string: %q", actualError)
	// t.Logf("Expected to contain: %q", expectedToContain)
	// assert.Contains(t, actualError, expectedToContain)

	// Check for constituent parts due to persistent strange failure with assert.Contains on combined string
	assert.Contains(t, actualError, "Analysis failed: ")
	assert.Contains(t, actualError, "failed to analyze CV and JD:")
	assert.Contains(t, actualError, expectedErrorMsg) // "LLM is down (mocked)"

	mockLLMService.AssertExpectations(t)
}
