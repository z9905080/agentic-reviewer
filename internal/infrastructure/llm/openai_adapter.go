package llm

import (
	"context"
	"cv-analyzer/internal/application/ports"
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"encoding/json" // For attempting to parse structured JSON from LLM
	"errors"
	"fmt"
	"os"
	"regexp" // For parsing score if not JSON
	"strconv" // For parsing score if not JSON
	"strings"
	"time"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
	"github.com/rs/zerolog/log" // For logging
)

type OpenAIAdapter struct {
	apiKey string
	client *openai.Client
}

func NewOpenAIAdapter(apiKey string) (*OpenAIAdapter, error) {
	l := log.With().Str("adapter", "OpenAIAdapter").Logger()
	effectiveAPIKey := apiKey
	if effectiveAPIKey == "" {
		l.Info().Msg("Direct API key not provided, checking environment variable OPENAI_API_KEY")
		effectiveAPIKey = os.Getenv("OPENAI_API_KEY")
	}

	if effectiveAPIKey == "" {
		l.Error().Msg("OpenAI API key is not configured (neither directly nor via ENV)")
		// This error message should match what NewOpenAIAdapter previously returned for consistency with tests if possible
		// Previous was: "OpenAI API key is not configured either directly or via OPENAI_API_KEY env var"
		// Let's keep it specific to this adapter's check.
		return nil, errors.New("OpenAI API key is not configured")
	}
	l.Info().Msg("OpenAIAdapter initialized with API key.")
	return &OpenAIAdapter{
		apiKey: effectiveAPIKey,
		client: openai.NewClient(effectiveAPIKey),
	}, nil
}

// Helper struct for attempting to parse a structured JSON response from LLM
type llmStructuredResponse struct {
	Score            float64  `json:"score"`
	Summary          string   `json:"summary"`
	MatchingKeywords []string `json:"matching_keywords"`
	MissingKeywords  []string `json:"missing_keywords"`
}

func (a *OpenAIAdapter) Analyze(cvData *cv.CV, jdData *jd.JobDescription) (*analysis.AnalysisResult, error) {
	l := log.With().Str("adapter", "OpenAIAdapter").Str("method", "Analyze").Logger()
	if cvData == nil || jdData == nil {
		l.Warn().Msg("CV or JD data is nil")
		return nil, errors.New("CV or JD data cannot be nil")
	}
	// Check for TextContent was previously here. Re-adding for robustness.
	if cvData.TextContent == "" {
        l.Warn().Str("cv_id", cvData.ID).Msg("CV text content is empty")
        return nil, errors.New("CV text content is empty")
    }
    if jdData.TextContent == "" {
        l.Warn().Str("jd_id", jdData.ID).Msg("JD text content is empty")
        return nil, errors.New("JD text content is empty")
    }

	l.Info().Str("cv_id", cvData.ID).Str("jd_id", jdData.ID).Msg("Starting analysis with OpenAI")

	// Constructing the prompt
	prompt := fmt.Sprintf(`Analyze the following CV based on the Job Description.
Provide a score from 0 to 100 indicating how well the CV matches the JD.
Provide a brief summary of the analysis.
List some matching keywords found in the CV relevant to the JD.
List some missing keywords in the CV that are important for the JD.

Format your response as a JSON object with the following keys: "score", "summary", "matching_keywords", "missing_keywords".

Job Description:
---
%s
---

CV:
---
%s
---
`, jdData.TextContent, cvData.TextContent)

	l.Debug().Str("model", openai.GPT3Dot5Turbo).Msg("Sending request to OpenAI ChatCompletion API")
	resp, err := a.client.CreateChatCompletion(
		context.Background(), // Ensure context is imported
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			// Temperature: 0.2,
			// MaxTokens:   500,
		},
	)

	if err != nil {
		l.Error().Err(err).Msg("OpenAI API request failed")
		return nil, fmt.Errorf("OpenAI API request failed: %w", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		l.Error().Msg("OpenAI returned no choices or empty message content")
		return nil, errors.New("no response choices from OpenAI")
	}

	llmOutput := resp.Choices[0].Message.Content
	l.Info().Int("llm_output_length", len(llmOutput)).Msg("Received response from OpenAI")
	l.Debug().Str("llm_raw_output", llmOutput).Msg("OpenAI Raw Output")

	var structuredResp llmStructuredResponse
	// Initialize AnalysisResult with some defaults and ensure all fields from domain are present
	analysisResult := analysis.AnalysisResult{
		ID:             uuid.NewString(),
		CV:             cvData, // Associate CV
		JobDescription: jdData, // Associate JD
		Score:          0.0,    // Default score
		FullReport:     llmOutput, // Always store the full raw output
		AnalyzedAt:     time.Now().UTC(),
	}

	jsonRegex := regexp.MustCompile("(?s)\\`\\`\\`json\\n(.*?)\\`\\`\\`")
	matches := jsonRegex.FindStringSubmatch(llmOutput)
	jsonString := llmOutput
	if len(matches) > 1 {
		jsonString = matches[1]
		l.Info().Msg("Extracted JSON block from LLM output")
	} else {
        if !strings.HasPrefix(strings.TrimSpace(jsonString), "{") {
            l.Warn().Msg("LLM output does not appear to be a JSON object, direct parsing might fail.")
        }
    }

	err = json.Unmarshal([]byte(jsonString), &structuredResp)
	if err != nil {
		l.Warn().Err(err).Str("json_string_attempted", jsonString).Msg("Failed to parse LLM output as JSON. Attempting basic extraction from raw output.")
		analysisResult.Summary = "LLM analysis complete. Could not parse structured data. Raw output in FullReport."

		scoreRegex := regexp.MustCompile(`(?i)["']?score["']?\s*:\s*([0-9.]+)`)
		scoreMatches := scoreRegex.FindStringSubmatch(llmOutput) // Parse from original llmOutput
		if len(scoreMatches) > 1 {
			if s, convErr := strconv.ParseFloat(scoreMatches[1], 64); convErr == nil {
				analysisResult.Score = s
				l.Info().Float64("extracted_score_via_regex", s).Msg("Extracted score using regex")
			} else {
				l.Warn().Err(convErr).Str("score_match", scoreMatches[1]).Msg("Failed to convert extracted score to float")
			}
		} else {
			l.Warn().Msg("Score not found via regex in LLM output.")
		}
	} else {
		l.Info().Msg("Successfully parsed structured JSON response from LLM")
		analysisResult.Score = structuredResp.Score
		analysisResult.Summary = structuredResp.Summary
		analysisResult.MatchingKeywords = structuredResp.MatchingKeywords
		analysisResult.MissingKeywords = structuredResp.MissingKeywords
	}

	// The following fields are now set when analysisResult is initialized or from structuredResp
	// analysisResult.CVID = cvData.ID // This is analysisResult.CV.ID
	// analysisResult.CVPath = cvData.FilePath // analysisResult.CV.FilePath
	// analysisResult.JDID = jdData.ID // analysisResult.JobDescription.ID
	// analysisResult.JDPath = jdData.FilePath // analysisResult.JobDescription.FilePath

	l.Info().Str("analysis_id", analysisResult.ID).Float64("score", analysisResult.Score).Msg("OpenAI analysis processing complete")
	return &analysisResult, nil
}

// Compile-time check
var _ ports.AnalysisService = (*OpenAIAdapter)(nil)
