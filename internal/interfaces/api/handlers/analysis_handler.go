package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath" // 用於創建臨時文件路徑
	"strings"       // Added for strings.Contains in error handling & isValidExtension

	"cv-analyzer/internal/application/usecases" // 確保路徑正確
	"cv-analyzer/internal/interfaces/api/dto"   // 確保路徑正確
	// "cv-analyzer/internal/domain/analysis" // Not strictly needed if using AnalysisResponse DTO

	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // 用於生成唯一文件名
	"time" // For AnalyzedAt
)

// AnalysisHandler 處理與 CV 和 JD 分析相關的 API 請求
type AnalysisHandler struct {
	analyzeUseCase usecases.AnalyzeCVUseCase // Using concrete type as per existing structure
}

// NewAnalysisHandler 創建一個新的 AnalysisHandler
func NewAnalysisHandler(analyzeUseCase usecases.AnalyzeCVUseCase) *AnalysisHandler {
	// In a typical setup with interfaces:
	// func NewAnalysisHandler(useCase ports.AnalyzeUseCase) *AnalysisHandler {
	return &AnalysisHandler{
		analyzeUseCase: analyzeUseCase,
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
	cvFile, cvHeader, err := c.Request.FormFile("cv")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "CV file is required"})
		return
	}
	defer cvFile.Close()

	jdFile, jdHeader, err := c.Request.FormFile("jd")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "JD file is required"})
		return
	}
	defer jdFile.Close()

	cvExt := filepath.Ext(cvHeader.Filename)
	jdExt := filepath.Ext(jdHeader.Filename)
	if !isValidExtension(cvExt) || !isValidExtension(jdExt) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Unsupported file type. Only .txt and .pdf are allowed."})
		return
	}

	tempDir := os.TempDir()
	cvTempFilePath, err := saveTempFile(cvFile, cvHeader, tempDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fmt.Sprintf("Failed to save CV file: %s", err.Error())})
		return
	}
	defer os.Remove(cvTempFilePath)

	jdTempFilePath, err := saveTempFile(jdFile, jdHeader, tempDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fmt.Sprintf("Failed to save JD file: %s", err.Error())})
		return
	}
	defer os.Remove(jdTempFilePath)

	analysisResult, err := h.analyzeUseCase.Execute(cvTempFilePath, jdTempFilePath)
	if err != nil {
		if strings.Contains(err.Error(), "unsupported file type") || strings.Contains(err.Error(), "content is empty") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fmt.Sprintf("Analysis failed: %s", err.Error())})
		}
		return
	}

	response := dto.AnalysisResponse{
		AnalysisID:       analysisResult.ID,
		CVID:             analysisResult.CV.ID,
		CVPath:           analysisResult.CV.FilePath, // This will be the temp path; consider if original filename is needed
		JDID:             analysisResult.JobDescription.ID,
		JDPath:           analysisResult.JobDescription.FilePath, // Also a temp path
		Score:            analysisResult.Score,
		Summary:          analysisResult.Summary,
		MatchingKeywords: analysisResult.MatchingKeywords,
		MissingKeywords:  analysisResult.MissingKeywords,
		AnalyzedAt:       time.Now(), // Set analysis time
	}
	c.JSON(http.StatusOK, response)
}

func saveTempFile(file multipart.File, header *multipart.FileHeader, tempDir string) (string, error) {
	fileName := uuid.NewString() + filepath.Ext(header.Filename)
	filePath := filepath.Join(tempDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("could not create temp file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		os.Remove(filePath) // Attempt to clean up partially written file
		return "", fmt.Errorf("could not save file content: %w", err)
	}
	return filePath, nil
}

func isValidExtension(ext string) bool {
	switch strings.ToLower(ext) {
	case ".txt", ".pdf":
		return true
	default:
		return false
	}
}
