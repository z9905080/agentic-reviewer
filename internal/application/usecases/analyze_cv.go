package usecases

import (
	"cv-analyzer/internal/application/ports"
	"cv-analyzer/internal/domain/analysis"
	"fmt"
	"github.com/google/uuid" // For generating unique IDs
)

// AnalyzeCVUseCase handles the business logic for analyzing a CV against a Job Description.
type AnalyzeCVUseCase struct {
	parserRouter    ports.FileParserRouter // Changed from individual parsers
	analysisService ports.AnalysisService
	storageService  ports.StorageService
}

// NewAnalyzeCVUseCase creates a new AnalyzeCVUseCase instance.
func NewAnalyzeCVUseCase(
	parserRouter ports.FileParserRouter, // Changed from individual parsers
	analysisService ports.AnalysisService,
	storageService ports.StorageService,
) *AnalyzeCVUseCase {
	return &AnalyzeCVUseCase{
		parserRouter:    parserRouter,
		analysisService: analysisService,
		storageService:  storageService,
	}
}

// Execute performs the CV analysis.
func (uc *AnalyzeCVUseCase) Execute(cvFilePath string, jdFilePath string) (*analysis.AnalysisResult, error) {
	// 1. Get CV Parser and Parse CV
	cvFileParser, err := uc.parserRouter.GetParser(cvFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get parser for CV file %s: %w", cvFilePath, err)
	}
	cvData, err := cvFileParser.ParseCV(cvFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CV file %s: %w", cvFilePath, err)
	}
	// Assign a new ID if not present (parsers should ideally handle ID generation)
	if cvData.ID == "" { // Parsers (PlainTextParser, PdfParser) now generate IDs. This might be redundant.
		cvData.ID = uuid.NewString()
	}
	// cvData.FilePath = cvFilePath // FilePath is set by the parser. This is redundant.

	// 2. Get JD Parser and Parse Job Description
	jdFileParser, err := uc.parserRouter.GetParser(jdFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get parser for JD file %s: %w", jdFilePath, err)
	}
	jdData, err := jdFileParser.ParseJD(jdFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JD file %s: %w", jdFilePath, err)
	}
	if jdData.ID == "" { // Parsers now generate IDs. This might be redundant.
		jdData.ID = uuid.NewString()
	}
	// jdData.FilePath = jdFilePath // FilePath is set by the parser. This is redundant.

	// 3. Perform Analysis
	analysisResult, err := uc.analysisService.Analyze(cvData, jdData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze CV and JD: %w", err)
	}
	if analysisResult.ID == "" {
		analysisResult.ID = uuid.NewString()
	}

	// 4. Store all entities
	if err := uc.storageService.SaveCV(cvData); err != nil {
		// Log error but continue, as analysis is the primary result
		fmt.Printf("Warning: failed to save CV %s: %v\n", cvData.ID, err)
	}
	if err := uc.storageService.SaveJD(jdData); err != nil {
		// Log error but continue
		fmt.Printf("Warning: failed to save JD %s: %v\n", jdData.ID, err)
	}
	if err := uc.storageService.SaveAnalysisResult(analysisResult); err != nil {
		return nil, fmt.Errorf("failed to save analysis result: %w", err)
	}

	return analysisResult, nil
}
