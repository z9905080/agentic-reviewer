package llm

import (
	// "context" // Potentially needed for Claude SDK - Removed as unused for now
	"errors"
	"fmt"
	"os"

	// "github.com/anthropic-ai/claude-sdk-go" // 假設有一個 Claude Go SDK
	// 如果沒有官方 SDK，可能需要使用 net/http 直接調用 API

	"cv-analyzer/internal/application/ports" // 確保路徑正確
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"github.com/google/uuid"
)

// ClaudeAdapter 實作 ports.AnalysisService
type ClaudeAdapter struct {
	apiKey string
	// client *claude.Client // 如果使用 SDK
}

// NewClaudeAdapter 創建一個新的 ClaudeAdapter 實例
func NewClaudeAdapter(apiKey string) (*ClaudeAdapter, error) {
	if apiKey == "" {
		apiKey = os.Getenv("CLAUDE_API_KEY")
		if apiKey == "" {
			return nil, errors.New("Claude API key is not configured")
		}
	}
	// client := claude.NewClient(apiKey) // 如果使用 SDK
	return &ClaudeAdapter{
		apiKey: apiKey,
		// client: client,
	}, nil
}

// Analyze 使用 Claude API 分析 CV 和 JD
// 方法簽名需與 ports.AnalysisService 介面一致
func (a *ClaudeAdapter) Analyze(cvData *cv.CV, jdData *jd.JobDescription) (*analysis.AnalysisResult, error) {
	if cvData == nil || jdData == nil {
		return nil, errors.New("CV or JD data cannot be nil")
	}

	// 這裡只是佔位邏輯，需要替換為實際的 Claude API 調用
	// 實際的 prompt 工程會在這裡發生
	// Now, we use the TextContent field populated by the parser.

	if cvData.TextContent == "" {
		return nil, errors.New("CV text content is empty")
	}
	if jdData.TextContent == "" {
		return nil, errors.New("JD text content is empty")
	}

	// Use the TextContent directly (example for prompt construction)
	// cvText := cvData.TextContent
	// jdText := jdData.TextContent

	// 範例：模擬 API 調用 (如果使用 SDK)
	// request := claude.NewCompletionRequest(claude.ModelLatest).
	// SetPrompt(fmt.Sprintf("Analyze this CV: %s based on this JD: %s", cvText, jdText)).
	// SetMaxTokensToSample(1000) // 示例參數
	//
	// resp, err := a.client.Complete(context.Background(), request)
	// if err != nil {
	// return nil, fmt.Errorf("Claude API request failed: %w", err)
	// }
	//
	// resultText := resp.Completion // 示例

	// 暫時返回一個模擬結果
	// Using cvData.TextContent and jdData.TextContent in the mock summary for demonstration
	resultText := fmt.Sprintf("Mock Claude Analysis for CV (%s) and JD (%s). They seem well-matched.", cvData.TextContent[:min(100, len(cvData.TextContent))], jdData.TextContent[:min(100, len(jdData.TextContent))])

	mockResult := &analysis.AnalysisResult{
		ID:             uuid.NewString(),
		CV:             cvData,
		JobDescription: jdData,
		Score:          0.855,
		Summary:        "This is a mock summary from Claude. " + resultText,
		MissingKeywords:  []string{"claude_mock_missing"},
		MatchingKeywords: []string{"claude_mock_matching1", "claude_mock_matching2"},
	}
	return mockResult, nil
}

// Helper function to prevent panic with slicing, useful for previews.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Compile-time check to ensure ClaudeAdapter implements the interface
var _ ports.AnalysisService = (*ClaudeAdapter)(nil)
