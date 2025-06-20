package llm

import (
	"context"
	"errors"
	"os"

	// 考慮使用官方或社群的 OpenAI Go SDK
	// "github.com/sashabaranov/go-openai" 是一個流行的選擇
	// 如果使用它，需要在 go.mod 中添加依賴
	openai "github.com/sashabaranov/go-openai"

	"cv-analyzer/internal/application/ports" // 確保路徑正確
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"fmt"
	"github.com/google/uuid"
)

// OpenAIAdapter 實作 ports.AnalysisService
type OpenAIAdapter struct {
	apiKey string // Store it if needed for other purposes, though client uses it
	client *openai.Client
}

// NewOpenAIAdapter 創建一個新的 OpenAIAdapter 實例
func NewOpenAIAdapter(apiKey string) (*OpenAIAdapter, error) {
	if apiKey == "" {
		// 實際應用中，API Key 應該從環境變數或設定檔讀取
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, errors.New("OpenAI API key is not configured either directly or via OPENAI_API_KEY env var")
		}
	}
	client := openai.NewClient(apiKey)
	return &OpenAIAdapter{
		apiKey: apiKey,
		client: client,
	}, nil
}

// Analyze 使用 OpenAI API 分析 CV 和 JD.
// This method signature must match the ports.AnalysisService interface.
func (a *OpenAIAdapter) Analyze(cvData *cv.CV, jdData *jd.JobDescription) (*analysis.AnalysisResult, error) {
	// 這裡只是佔位邏輯，需要替換為實際的 OpenAI API 調用
	// 實際的 prompt 工程會在這裡發生
	// It would involve extracting text from cvData and jdData.
	// Now, we use the TextContent field populated by the parser.

	if cvData == nil || jdData == nil {
		return nil, errors.New("CV or JD data cannot be nil")
	}
	if cvData.TextContent == "" {
		return nil, errors.New("CV text content is empty")
	}
	if jdData.TextContent == "" {
		return nil, errors.New("JD text content is empty")
	}

	// Use the TextContent directly
	cvText := cvData.TextContent
	jdText := jdData.TextContent

	// Simulate an OpenAI API call and response.
	// In a real scenario, you would uncomment and adapt the following:
	/*
	prompt := fmt.Sprintf("Analyze the following CV based on the Job Description. "+
		"Provide a matching score (0.0-1.0), a brief summary of the match, "+
		"a list of matching keywords, and a list of missing keywords from the CV that are present in the JD."+
		"\n\nCV Content:\n%s\n\nJob Description Content:\n%s", cvText, jdText)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo, // Or any other model like GPT4TurboPreview
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		// MaxTokens: 1000, // Optional: Adjust as needed
		// Temperature: 0.7, // Optional: Adjust for creativity vs. determinism
	}

	resp, err := a.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API request failed: %w", err)
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return nil, errors.New("no response choices or empty content from OpenAI")
	}
	analysisContent := resp.Choices[0].Message.Content

	// Here you would parse 'analysisContent' to extract score, summary, matching/missing keywords.
	// This parsing logic can be complex depending on the LLM's output format.
	// For now, we continue with mock data.
	*/
	mockAnalysisContent := fmt.Sprintf("Mock OpenAI Analysis for CV %s (%s) and JD %s (%s). They are a good match.", cvData.ID, cvText, jdData.ID, jdText)
	mockScore := 0.85

	// Construct and return an AnalysisResult
	result := &analysis.AnalysisResult{
		ID:               uuid.NewString(), // Generate a new ID for the analysis
		CV:               cvData,
		JobDescription:   jdData,
		Score:            mockScore,
		Summary:          mockAnalysisContent,
		MissingKeywords:  []string{"example_missing_from_openai_adapter"},
		MatchingKeywords: []string{"example_matching_from_openai_adapter"},
	}

	return result, nil
}

// Compile-time check to ensure OpenAIAdapter implements the AnalysisService interface.
var _ ports.AnalysisService = (*OpenAIAdapter)(nil)
