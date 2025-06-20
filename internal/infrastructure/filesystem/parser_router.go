package filesystem

import (
	"path/filepath" // 用於獲取文件擴展名
	"strings"
	"fmt"

	"cv-analyzer/internal/application/ports"
)

// ParserRouter 根據文件類型選擇合適的文件解析器
type ParserRouter struct {
	plainTextParser *PlainTextParser
	pdfParser       *PdfParser
	// docxParser *DocxParser // 未來可以添加
}

// NewParserRouter 創建一個新的 ParserRouter
// Refactored to accept parser dependencies for proper DI with Wire.
func NewParserRouter(plainParser *PlainTextParser, pdfParser *PdfParser) (*ParserRouter, error) {
	if plainParser == nil {
		return nil, fmt.Errorf("PlainTextParser cannot be nil for ParserRouter")
	}
	if pdfParser == nil {
		return nil, fmt.Errorf("PdfParser cannot be nil for ParserRouter")
	}
	return &ParserRouter{
		plainTextParser: plainParser,
		pdfParser:       pdfParser,
	}, nil
}

// GetParser 根據文件路徑（擴展名）返回合適的 FileParser
// 注意：返回的是 ports.FileParser 介面類型
func (r *ParserRouter) GetParser(filePath string) (ports.FileParser, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt":
		return r.plainTextParser, nil
	case ".pdf":
		return r.pdfParser, nil
	// case ".docx":
	//  if r.docxParser == nil {
	//      return nil, fmt.Errorf("DOCX parser not initialized")
	//  }
	// 	return r.docxParser, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s. Only .txt and .pdf are supported", ext)
	}
}

// Compile-time check to ensure ParserRouter implements the FileParserRouter interface.
var _ ports.FileParserRouter = (*ParserRouter)(nil)
