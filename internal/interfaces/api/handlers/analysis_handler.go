package handlers

import (
	"cv-analyzer/internal/application/usecases"
	"cv-analyzer/internal/interfaces/api/dto"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time" // For response DTO

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log" // Using package-level logger
)

type AnalysisHandler struct {
	analyzeUseCase usecases.AnalyzeCVUseCase // Now this should be the interface type
}

func NewAnalysisHandler(analyzeUseCase usecases.AnalyzeCVUseCase) *AnalysisHandler { // Parameter is interface
	return &AnalysisHandler{
		analyzeUseCase: analyzeUseCase, // Assigning interface to interface field
	}
}

// HandleAnalyzeCV godoc
// @Summary Analyze CV and JD by uploading files
// @Description Uploads a CV and a Job Description file (PDF or TXT) for analysis, processing, and scoring.
// @Tags analysis
// @Accept multipart/form-data
// @Produce json
// @Param cv formData file true "CV file (PDF or TXT)"
// @Param jd formData file true "Job Description file (PDF or TXT)"
// @Success 200 {object} dto.AnalysisResponse "Successful analysis"
// @Failure 400 {object} dto.ErrorResponse "Bad request (e.g., missing files, unsupported file type, invalid file)"
// @Failure 500 {object} dto.ErrorResponse "Internal server error (e.g., failed to save file, analysis process failed)"
// @Router /api/v1/analysis/cv [post]
func (h *AnalysisHandler) HandleAnalyzeCV(c *gin.Context) {
	// Create a logger instance with initial context for this request
	// Request ID could be added from a middleware if available: c.GetString("requestID")
	l := log.With().Str("handler", "AnalysisHandler").Str("method", "HandleAnalyzeCV").Logger()
	l.Info().Msg("Received new CV analysis request")

	cvFile, cvHeader, err := c.Request.FormFile("cv")
	if err != nil {
		l.Warn().Err(err).Msg("CV file is required")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "CV file is required"})
		return
	}
	defer cvFile.Close()
	// Update logger context with CV filename
	l = l.With().Str("cv_filename", cvHeader.Filename).Logger()

	jdFile, jdHeader, err := c.Request.FormFile("jd")
	if err != nil {
		l.Warn().Err(err).Msg("JD file is required")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "JD file is required"})
		return
	}
	defer jdFile.Close()
	// Update logger context with JD filename
	l = l.With().Str("jd_filename", jdHeader.Filename).Logger()
	l.Info().Msg("CV and JD files received")


	cvExt := filepath.Ext(cvHeader.Filename)
	jdExt := filepath.Ext(jdHeader.Filename)
	if !isValidExtension(cvExt) || !isValidExtension(jdExt) {
		errMsg := "Unsupported file type. Only .txt and .pdf are allowed."
		l.Warn().Str("cv_ext", cvExt).Str("jd_ext", jdExt).Msg(errMsg)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: errMsg})
		return
	}

	tempDir := os.TempDir()
	cvTempFilePath, err := saveTempFile(cvFile, cvHeader, tempDir)
	if err != nil {
		l.Error().Err(err).Str("original_filename", cvHeader.Filename).Msg("Failed to save CV temp file")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fmt.Sprintf("Failed to save CV file: %s", err.Error())})
		return
	}
	defer os.Remove(cvTempFilePath)
	l.Debug().Str("cv_temp_path", cvTempFilePath).Msg("CV file saved to temporary location")

	jdTempFilePath, err := saveTempFile(jdFile, jdHeader, tempDir)
	if err != nil {
		l.Error().Err(err).Str("original_filename", jdHeader.Filename).Msg("Failed to save JD temp file")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fmt.Sprintf("Failed to save JD file: %s", err.Error())})
		return
	}
	defer os.Remove(jdTempFilePath)
	l.Debug().Str("jd_temp_path", jdTempFilePath).Msg("JD file saved to temporary location")

	l.Info().Msg("Calling analysis use case")
	analysisResult, err := h.analyzeUseCase.Execute(cvTempFilePath, jdTempFilePath)
	if err != nil {
		l.Error().Err(err).Msg("Analysis use case execution failed")
		// Check if the error is from parser router (unsupported file type) or other parsing errors
		if strings.Contains(err.Error(), "unsupported file type") ||
		   strings.Contains(err.Error(), "error parsing CV") ||
		   strings.Contains(err.Error(), "error parsing JD") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		} else { // Assume other errors are internal server errors (e.g. LLM failure, DB failure)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fmt.Sprintf("Analysis failed: %s", err.Error())})
		}
		return
	}

	// Correctly map from analysis.AnalysisResult (which has nested CV & JD) to dto.AnalysisResponse
	response := dto.AnalysisResponse{
		AnalysisID:       analysisResult.ID,
		Score:            analysisResult.Score,
		Summary:          analysisResult.Summary,
		MissingKeywords:  analysisResult.MissingKeywords,
		MatchingKeywords: analysisResult.MatchingKeywords,
		CVID:             analysisResult.CV.ID,
		CVPath:           analysisResult.CV.FilePath, // This is the temp path
		JDID:             analysisResult.JobDescription.ID,
		JDPath:           analysisResult.JobDescription.FilePath, // This is the temp path
		AnalyzedAt:       time.Now().UTC(), // Or use analysisResult.AnalyzedAt if set by use case
	}
	l.Info().Str("analysis_id", response.AnalysisID).Msg("Analysis successful")
	c.JSON(http.StatusOK, response)
}

// saveTempFile helper remains the same as in subtask 9
func saveTempFile(file multipart.File, header *multipart.FileHeader, tempDir string) (string, error) {
	fileName := uuid.NewString() + filepath.Ext(header.Filename)
	filePath := filepath.Join(tempDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("could not create temp file %s: %w", filePath, err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("could not save file content to %s: %w", filePath, err)
	}
	return filePath, nil
}

// isValidExtension helper remains the same
func isValidExtension(ext string) bool {
	switch strings.ToLower(ext) {
	case ".txt", ".pdf":
		return true
	default:
		return false
	}
}
