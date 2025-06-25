package usecases

import (
	"cv-analyzer/internal/application/ports"
	"cv-analyzer/internal/domain/analysis"
	// "cv-analyzer/internal/domain/cv" // Not needed if cvModel is used directly
	// "cv-analyzer/internal/domain/jd" // Not needed if jdModel is used directly
	"fmt"
	// "time" // Removed as AnalyzedAt is set in handler

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// AnalyzeCVUseCase defines the interface for the use case.
// This was previously a struct type directly. Defining an interface is good practice.
type AnalyzeCVUseCase interface {
	Execute(cvFilePath string, jdFilePath string) (*analysis.AnalysisResult, error)
}

// analyzeCVUseCase implements AnalyzeCVUseCase.
// Renamed from AnalyzeCVUseCase (struct) to analyzeCVUseCase (struct)
type analyzeCVUseCase struct {
	parserRouter    ports.FileParserRouter
	analysisService ports.AnalysisService
	storageService  ports.StorageService
}

// NewAnalyzeCVUseCase creates a new AnalyzeCVUseCase instance.
// Returns the interface type.
func NewAnalyzeCVUseCase(
	parserRouter ports.FileParserRouter,
	analysisService ports.AnalysisService,
	storageService ports.StorageService,
) AnalyzeCVUseCase { // Return interface type
	return &analyzeCVUseCase{ // Return pointer to struct
		parserRouter:    parserRouter,
		analysisService: analysisService,
		storageService:  storageService,
	}
}

func (uc *analyzeCVUseCase) Execute(cvFilePath string, jdFilePath string) (*analysis.AnalysisResult, error) {
	l := log.With().Str("usecase", "AnalyzeCV").Str("cv_path", cvFilePath).Str("jd_path", jdFilePath).Logger()
	l.Info().Msg("Starting CV analysis")

	cvParser, err := uc.parserRouter.GetParser(cvFilePath)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get CV parser")
		return nil, fmt.Errorf("error getting parser for CV file %s: %w", cvFilePath, err)
	}
	jdParser, err := uc.parserRouter.GetParser(jdFilePath)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get JD parser")
		return nil, fmt.Errorf("error getting parser for JD file %s: %w", jdFilePath, err)
	}

	l.Debug().Msg("Parsing CV")
	cvModel, err := cvParser.ParseCV(cvFilePath)
	if err != nil {
		l.Error().Err(err).Msg("Failed to parse CV")
		return nil, fmt.Errorf("error parsing CV file %s: %w", cvFilePath, err)
	}
	l.Info().Str("cv_id", cvModel.ID).Msg("CV parsed successfully")

	l.Debug().Msg("Parsing JD")
	jdModel, err := jdParser.ParseJD(jdFilePath)
	if err != nil {
		l.Error().Err(err).Msg("Failed to parse JD")
		return nil, fmt.Errorf("error parsing JD file %s: %w", jdFilePath, err)
	}
	l.Info().Str("jd_id", jdModel.ID).Msg("JD parsed successfully")

	// Save CV and JD metadata (as per original use case logic)
	if err := uc.storageService.SaveCV(cvModel); err != nil {
		l.Warn().Err(err).Str("cv_id", cvModel.ID).Msg("Failed to save CV metadata, continuing analysis")
	}
	if err := uc.storageService.SaveJD(jdModel); err != nil {
		l.Warn().Err(err).Str("jd_id", jdModel.ID).Msg("Failed to save JD metadata, continuing analysis")
	}

	l.Debug().Msg("Performing LLM analysis")
	analysisResult, err := uc.analysisService.Analyze(cvModel, jdModel)
	if err != nil {
		l.Error().Err(err).Msg("LLM analysis failed")
		return nil, fmt.Errorf("failed to analyze CV and JD: %w", err)
	}

	// Populate/ensure fields in AnalysisResult
	// The analysisService might already populate some of these.
	// If analysisResult.ID is empty, generate one.
	if analysisResult.ID == "" {
	    analysisResult.ID = uuid.NewString()
    }
	// Ensure CV and JD are associated if not already done by the LLM service.
	analysisResult.CV = cvModel
	analysisResult.JobDescription = jdModel
	// analysisResult.AnalyzedAt = time.Now().UTC() // This is now set in the handler for the DTO

	l.Info().Str("analysis_id", analysisResult.ID).Msg("LLM analysis complete")

	l.Debug().Str("analysis_id", analysisResult.ID).Msg("Saving analysis result")
	if err := uc.storageService.SaveAnalysisResult(analysisResult); err != nil {
		l.Error().Err(err).Str("analysis_id", analysisResult.ID).Msg("Failed to save analysis result")
		return nil, fmt.Errorf("failed to save analysis result: %w", err)
	}
	l.Info().Str("analysis_id", analysisResult.ID).Msg("Analysis result saved")

	return analysisResult, nil
}
