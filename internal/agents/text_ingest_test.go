package agents

import (
	"path/filepath"
	"testing"
)

func TestTextIngestAgent_ProcessTxtFile(t *testing.T) {
	agent := NewTextIngestAgent()

	// Test English text file
	englishFile := filepath.Join("..", "..", "testdata", "samples", "sample.txt")
	content, err := agent.ProcessFile(englishFile)
	if err != nil {
		t.Fatalf("Failed to process English text file: %v", err)
	}

	if len(content.Paragraphs) == 0 {
		t.Error("No paragraphs found in English text")
	}

	if content.WordCount == 0 {
		t.Error("No words counted in English text")
	}

	if content.Language != "en-US" {
		t.Errorf("Expected language en-US, got %s", content.Language)
	}

	t.Logf("English file: %d paragraphs, %d words, language: %s",
		len(content.Paragraphs), content.WordCount, content.Language)
}

func TestTextIngestAgent_ProcessGreekFile(t *testing.T) {
	agent := NewTextIngestAgent()

	// Test Greek text file
	greekFile := filepath.Join("..", "..", "testdata", "samples", "greek.txt")
	content, err := agent.ProcessFile(greekFile)
	if err != nil {
		t.Fatalf("Failed to process Greek text file: %v", err)
	}

	if len(content.Paragraphs) == 0 {
		t.Error("No paragraphs found in Greek text")
	}

	if content.WordCount == 0 {
		t.Error("No words counted in Greek text")
	}

	if content.Language != "el-GR" {
		t.Errorf("Expected language el-GR, got %s", content.Language)
	}

	t.Logf("Greek file: %d paragraphs, %d words, language: %s",
		len(content.Paragraphs), content.WordCount, content.Language)
}

func TestTextIngestAgent_ValidateContent(t *testing.T) {
	agent := NewTextIngestAgent()

	// Test valid content
	validContent := &TextContent{
		Paragraphs: []string{"Hello world", "This is a test"},
		WordCount:  5,
		Language:   "en-US",
	}

	if err := agent.ValidateContent(validContent); err != nil {
		t.Errorf("Valid content failed validation: %v", err)
	}

	// Test invalid content (no paragraphs)
	invalidContent := &TextContent{
		Paragraphs: []string{},
		WordCount:  0,
		Language:   "en-US",
	}

	if err := agent.ValidateContent(invalidContent); err == nil {
		t.Error("Invalid content passed validation")
	}
}
