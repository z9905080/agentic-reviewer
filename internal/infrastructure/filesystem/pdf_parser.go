package filesystem

import (
	"bytes"
	"errors"
	"fmt"
	"os" // Added import
	"strings"

	"rsc.io/pdf"

	"cv-analyzer/internal/application/ports"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"github.com/google/uuid"
)

type PdfParser struct{}

func NewPdfParser() *PdfParser {
	return &PdfParser{}
}

func (p *PdfParser) parseContent(pdfData []byte) (string, error) {
	r, err := pdf.NewReader(bytes.NewReader(pdfData), int64(len(pdfData)))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var content strings.Builder
	numPages := r.NumPage()
	if numPages == 0 {
		return "", errors.New("PDF has no pages")
	}
	for i := 1; i <= numPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}

		texts := page.Content().Text
		for _, t := range texts {
			content.WriteString(t.S)
			// Consider adding spaces based on positioning if needed for readability
			// content.WriteString(" ")
		}
		if len(texts) > 0 && i < numPages { // Add a newline between pages if text was extracted
		    content.WriteString("\n")
		}
	}

	extractedStr := content.String()
	if strings.TrimSpace(extractedStr) == "" {
		// It's possible a PDF contains no text (e.g., scanned images)
		// Decide if this is an error or just an empty result.
		// For now, let's return empty string and no error if parsing succeeded but found no text.
		// return "", errors.New("no text content found in PDF")
	}
	return extractedStr, nil
}

func (p *PdfParser) ParseCV(filePath string) (*cv.CV, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file %s: %w", filePath, err)
	}
	textContent, err := p.parseContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from PDF %s: %w", filePath, err)
	}
	newCv := &cv.CV{
		ID:          uuid.NewString(),
		FilePath:    filePath,
		TextContent: textContent,
	}
	return newCv, nil
}

func (p *PdfParser) ParseJD(filePath string) (*jd.JobDescription, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file %s: %w", filePath, err)
	}
	textContent, err := p.parseContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from PDF %s: %w", filePath, err)
	}
	newJd := &jd.JobDescription{
		ID:          uuid.NewString(),
		FilePath:    filePath,
		TextContent: textContent,
	}
	return newJd, nil
}

var _ ports.CVParser = (*PdfParser)(nil)
var _ ports.JDParser = (*PdfParser)(nil)
