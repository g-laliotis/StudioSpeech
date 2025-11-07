package agents

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
	"github.com/unidoc/unioffice/document"
)

// TextContent represents processed text content
type TextContent struct {
	Paragraphs []string
	Language   string // Detected or specified language
	WordCount  int
	Source     string // Original file path
}

// TextIngestAgent handles reading and processing text files
type TextIngestAgent struct{}

// NewTextIngestAgent creates a new text ingestion agent
func NewTextIngestAgent() *TextIngestAgent {
	return &TextIngestAgent{}
}

// ProcessFile reads and processes a text file (.txt, .docx, or .pdf)
func (t *TextIngestAgent) ProcessFile(filePath string) (*TextContent, error) {
	// Validate file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}
	
	ext := strings.ToLower(filepath.Ext(filePath))
	
	var content *TextContent
	var err error
	
	switch ext {
	case ".txt":
		content, err = t.processTxtFile(filePath)
	case ".docx":
		content, err = t.processDocxFile(filePath)
	case ".pdf":
		content, err = t.processPdfFile(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s (supported: .txt, .docx, .pdf)", ext)
	}
	
	if err != nil {
		return nil, err
	}
	
	content.Source = filePath
	content.WordCount = t.countWords(content.Paragraphs)
	
	return content, nil
}

// processTxtFile reads and processes a plain text file
func (t *TextIngestAgent) processTxtFile(filePath string) (*TextContent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open text file: %w", err)
	}
	defer file.Close()
	
	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read text file: %w", err)
	}
	
	// Validate UTF-8 encoding
	if !utf8.Valid(content) {
		return nil, fmt.Errorf("file is not valid UTF-8 encoded")
	}
	
	text := string(content)
	paragraphs := t.splitIntoParagraphs(text)
	
	return &TextContent{
		Paragraphs: paragraphs,
		Language:   t.detectLanguage(text),
	}, nil
}

// processDocxFile reads and processes a Microsoft Word document
func (t *TextIngestAgent) processDocxFile(filePath string) (*TextContent, error) {
	doc, err := document.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DOCX file: %w", err)
	}
	defer doc.Close()
	
	var paragraphs []string
	var allText strings.Builder
	
	// Extract text from all paragraphs
	for _, para := range doc.Paragraphs() {
		// Extract text from paragraph runs
		var paraText strings.Builder
		for _, run := range para.Runs() {
			paraText.WriteString(run.Text())
		}
		
		text := strings.TrimSpace(paraText.String())
		if text != "" {
			paragraphs = append(paragraphs, text)
			allText.WriteString(text + " ")
		}
	}
	
	if len(paragraphs) == 0 {
		return nil, fmt.Errorf("no text content found in DOCX file")
	}
	
	return &TextContent{
		Paragraphs: paragraphs,
		Language:   t.detectLanguage(allText.String()),
	}, nil
}

// processPdfFile reads and processes a PDF document
func (t *TextIngestAgent) processPdfFile(filePath string) (*TextContent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer file.Close()
	
	// Get file info for PDF reader
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get PDF file info: %w", err)
	}
	
	// Read PDF content
	reader, err := pdf.NewReader(file, fileInfo.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}
	
	var paragraphs []string
	var allText strings.Builder
	
	// Extract text from all pages
	for i := 1; i <= reader.NumPage(); i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		
		// Get page content
		pageText, err := page.GetPlainText(nil)
		if err != nil {
			continue // Skip pages with extraction errors
		}
		
		// Clean and process page text
		pageText = strings.TrimSpace(pageText)
		if pageText != "" {
			// Split page into paragraphs
			pageParagraphs := t.splitIntoParagraphs(pageText)
			paragraphs = append(paragraphs, pageParagraphs...)
			allText.WriteString(pageText + " ")
		}
	}
	
	if len(paragraphs) == 0 {
		return nil, fmt.Errorf("no text content found in PDF file")
	}
	
	return &TextContent{
		Paragraphs: paragraphs,
		Language:   t.detectLanguage(allText.String()),
	}, nil
}

// splitIntoParagraphs splits text into paragraphs, preserving structure
func (t *TextIngestAgent) splitIntoParagraphs(text string) []string {
	var paragraphs []string
	
	// Split by double newlines (paragraph breaks)
	parts := strings.Split(text, "\n\n")
	
	for _, part := range parts {
		// Clean up the paragraph
		para := strings.TrimSpace(part)
		// Replace single newlines with spaces (soft breaks)
		para = strings.ReplaceAll(para, "\n", " ")
		// Normalize multiple spaces
		para = strings.Join(strings.Fields(para), " ")
		
		if para != "" {
			paragraphs = append(paragraphs, para)
		}
	}
	
	// If no paragraph breaks found, split by single newlines
	if len(paragraphs) <= 1 && strings.Contains(text, "\n") {
		scanner := bufio.NewScanner(strings.NewReader(text))
		paragraphs = nil // Reset
		
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				paragraphs = append(paragraphs, line)
			}
		}
	}
	
	return paragraphs
}

// detectLanguage performs simple heuristic language detection
func (t *TextIngestAgent) detectLanguage(text string) string {
	// Simple heuristic based on character patterns
	text = strings.ToLower(text)
	
	// Count Greek characters
	greekCount := 0
	englishCount := 0
	
	for _, r := range text {
		if r >= 'α' && r <= 'ω' || r >= 'Α' && r <= 'Ω' {
			greekCount++
		} else if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			englishCount++
		}
	}
	
	// Determine language based on character distribution
	if greekCount > englishCount {
		return "el-GR"
	} else if englishCount > 0 {
		return "en-US"
	}
	
	return "auto" // Unknown
}

// countWords counts the total number of words in all paragraphs
func (t *TextIngestAgent) countWords(paragraphs []string) int {
	totalWords := 0
	for _, para := range paragraphs {
		words := strings.Fields(para)
		totalWords += len(words)
	}
	return totalWords
}

// ValidateContent performs basic content validation
func (t *TextIngestAgent) ValidateContent(content *TextContent) error {
	if content == nil {
		return fmt.Errorf("content is nil")
	}
	
	if len(content.Paragraphs) == 0 {
		return fmt.Errorf("no paragraphs found")
	}
	
	if content.WordCount == 0 {
		return fmt.Errorf("no words found")
	}
	
	// Check for reasonable content length
	if content.WordCount > 50000 {
		return fmt.Errorf("content too long: %d words (max 50,000)", content.WordCount)
	}
	
	return nil
}