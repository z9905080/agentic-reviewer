//go:build wireinject
// +build wireinject

package main

import (
	"cv-analyzer/internal/application/ports"
	"cv-analyzer/internal/application/usecases"
	"cv-analyzer/internal/infrastructure/filesystem"
	"cv-analyzer/internal/infrastructure/llm" // Assuming dummy or actual adapters are here
	"cv-analyzer/internal/infrastructure/persistence" // If we had a real storage service
	apiHandlers "cv-analyzer/internal/interfaces/api/handlers"
	"cv-analyzer/internal/interfaces/api/routes" // For router setup
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"os" // For reading env vars for API keys potentially
	"fmt" // For provideDummyAnalysisService panic
)

// --- Providers for Infrastructure Layer ---

// func providePlainTextParser() *filesystem.PlainTextParser {
// 	return filesystem.NewPlainTextParser()
// }

// func providePdfParser() *filesystem.PdfParser {
// 	return filesystem.NewPdfParser()
// }

// Provider for ParserRouter
// func provideParserRouter(plainParser *filesystem.PlainTextParser, pdfParser *filesystem.PdfParser) (ports.FileParserRouter, error) {
//     // NewParserRouter in the subtask for step 5 created new instances.
//     // For wire, it's better if NewParserRouter accepts its dependencies.
//     // Let's assume NewParserRouter is refactored or we construct it here.
//     // For now, using a simplified constructor that matches previous setup:
// 	return filesystem.NewParserRouter() // This internally news up parsers, not ideal for wire.
//                                         // We'll address this if wire complains or in a follow-up.
//                                         // IDEAL: filesystem.NewParserRouter(plainParser, pdfParser)
// }


// Provider for LLM Analysis Service
// Option 1: Dummy service
func provideDummyAnalysisService() ports.AnalysisService {
    // Assuming AnalysisServiceDummy exists and has a constructor
    // return llm.NewAnalysisServiceDummy()
    // For now, let's assume we want to wire up OpenAI by default if API key is present
    apiKey := os.Getenv("OPENAI_API_KEY")
    // if apiKey == "" { // NewOpenAIAdapter now handles empty key by checking env itself.
    //     // Fallback to dummy if no key, or handle error
    //     // This logic might be better outside the provider or in a config struct
    //     // return llm.NewDummyAnalysisService() // Placeholder from previous setup
    // }

    oa, err := llm.NewOpenAIAdapter(apiKey)
    if err != nil {
        // This is tricky in a provider. Panicking or returning a default/error.
        // For wire, providers should ideally not fail or handle config this way.
        // A Config struct is better.
        // For now, let's assume it can return an error.
        // return nil, err // Wire will handle this if the injector returns error
        // Forcing a panic if critical service fails and no fallback configured
        panic(fmt.Sprintf("Failed to init OpenAIAdapter: %v. Make sure OPENAI_API_KEY is set if intending to use it, or configure a fallback dummy service.", err))
    }
    return oa
}

// --- Providers for Application Layer ---

func provideAnalyzeCVUseCase(
	parserRouter ports.FileParserRouter,
	analysisService ports.AnalysisService,
	storageService ports.StorageService, // Now a proper dependency
) *usecases.AnalyzeCVUseCase { // Return pointer type for consistency
	return usecases.NewAnalyzeCVUseCase(parserRouter, analysisService, storageService)
}


// --- Providers for Interface Layer (API Handlers & Router) ---

func provideAnalysisHandler(useCase *usecases.AnalyzeCVUseCase) *apiHandlers.AnalysisHandler { // UseCase should be pointer
	return apiHandlers.NewAnalysisHandler(*useCase) // NewAnalysisHandler expects a concrete usecases.AnalyzeCVUseCase
}

// provideGinEngine initializes the Gin engine and sets up routes.
// Now depends on AnalysisHandler, which is provided by Wire.
func provideGinEngine(analysisHandler *apiHandlers.AnalysisHandler) *gin.Engine {
	// routes.SetupRouter now takes *apiHandlers.AnalysisHandler
	return routes.SetupRouter(analysisHandler)
}


// --- Main Injector ---
// initializeAPI is the main injector function that Wire will implement.
func initializeAPI() (*gin.Engine, error) {
	wire.Build(
		// Filesystem Parsers
        filesystem.NewPlainTextParser,
        filesystem.NewPdfParser,

        // ParserRouter - Assuming NewParserRouter needs to be refactored to accept dependencies
        // For now, the script had this:
        // wire.NewSet(
        //     filesystem.NewParserRouter, // This is func() (*ParserRouter, error)
        //     wire.Bind(new(ports.FileParserRouter), new(*filesystem.ParserRouter)),
        // ),
        // This is correct if NewParserRouter doesn't take plain/pdf parsers as args.
        // The current NewParserRouter in filesystem/parser_router.go is:
        // func NewParserRouter() (*ParserRouter, error) { /* internally news up PlainTextParser and PdfParser */ }
        // This is an anti-pattern for DI.
        // Let's assume it's refactored: func NewParserRouter(p *PlainTextParser, pdf *PdfParser) (*ParserRouter, error)
        // If so, the providers are:
        // filesystem.NewPlainTextParser, (already listed)
        // filesystem.NewPdfParser, (already listed)
        // filesystem.NewParserRouter, // Now func(p *PlainTextParser, pdf *PdfParser) (*ParserRouter, error)
        // wire.Bind(new(ports.FileParserRouter), new(*filesystem.ParserRouter)),
        // This will be an issue if NewParserRouter is not refactored.
        // filesystem.NewPlainTextParser, // Already listed globally or should be.
        // filesystem.NewPdfParser,       // Already listed globally or should be.
        filesystem.NewParserRouter,    // Wire will use NewPlainTextParser and NewPdfParser automatically.
        wire.Bind(new(ports.FileParserRouter), new(*filesystem.ParserRouter)),

		// LLM Service Provider
		provideDummyAnalysisService, // Returns ports.AnalysisService
        // wire.Bind(new(ports.AnalysisService), new(*llm.OpenAIAdapter)), // This would bind if provideDummyAnalysisService returned *llm.OpenAIAdapter

		// Persistence
		persistence.NewInMemoryStorage,
		wire.Bind(new(ports.StorageService), new(*persistence.InMemoryStorage)),

		provideAnalyzeCVUseCase,
		provideAnalysisHandler, // Now provideAnalysisHandler is needed by provideGinEngine.
		provideGinEngine,
	)
	return nil, nil // Wire will fill this in
}
