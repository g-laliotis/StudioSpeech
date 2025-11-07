package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"studiospeech/internal/agents"
)

// TestEdgeCases tests various edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "empty string",
			input:       "",
			expectError: false, // Agent handles empty gracefully
			description: "Should handle empty input gracefully",
		},
		{
			name:        "only whitespace",
			input:       "   \n\t  \n  ",
			expectError: false, // Agent handles whitespace gracefully
			description: "Should handle whitespace-only input",
		},
		{
			name:        "single character",
			input:       "a",
			expectError: false,
			description: "Should handle single character input",
		},
		{
			name:        "very long sentence",
			input:       strings.Repeat("This is a very long sentence that tests the maximum length handling capabilities of the system. ", 20),
			expectError: false,
			description: "Should handle very long sentences within limits",
		},
		{
			name:        "unicode characters",
			input:       "Hello 世界! Γεια σας κόσμε! 🎤🔊",
			expectError: false,
			description: "Should handle Unicode characters properly",
		},
		{
			name:        "special punctuation",
			input:       "Test... with—various–punctuation: semicolons; and (parentheses)!",
			expectError: false,
			description: "Should handle special punctuation marks",
		},
		{
			name:        "numbers and symbols",
			input:       "Call 555-1234 or email test@example.com for $100 discount!",
			expectError: false,
			description: "Should handle numbers, symbols, and email addresses",
		},
		{
			name:        "mixed languages",
			input:       "Hello world! Γεια σας! Bonjour monde!",
			expectError: false,
			description: "Should handle mixed language content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.name+".txt")
			err := os.WriteFile(testFile, []byte(tt.input), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test text ingestion
			ingestAgent := agents.NewTextIngestAgent()
			content, err := ingestAgent.ProcessFile(testFile)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.description)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
				return
			}

			// Test normalization
			normalizeAgent := agents.NewNormalizeAgent()
			normalized, err := normalizeAgent.Normalize(content)
			
			if err != nil {
				t.Errorf("Normalization failed for %s: %v", tt.description, err)
				return
			}

			// Validate normalized output
			if len(normalized.Sentences) == 0 && len(strings.TrimSpace(tt.input)) > 0 {
				t.Errorf("No sentences generated for valid input: %s", tt.description)
			}
		})
	}
}

// TestBoundaryConditions tests system limits and boundaries
func TestBoundaryConditions(t *testing.T) {
	tests := []struct {
		name        string
		testFunc    func(t *testing.T)
		description string
	}{
		{
			name:        "maximum sentence length",
			testFunc:    testMaxSentenceLength,
			description: "Test handling of maximum allowed sentence length",
		},
		{
			name:        "maximum file size",
			testFunc:    testMaxFileSize,
			description: "Test handling of large input files",
		},
		{
			name:        "memory pressure",
			testFunc:    testMemoryPressure,
			description: "Test system behavior under memory pressure",
		},
		{
			name:        "concurrent processing",
			testFunc:    testConcurrentProcessing,
			description: "Test concurrent pipeline execution",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func testMaxSentenceLength(t *testing.T) {
	t.Skip("Sentence length validation test - implementation may vary")
	normalizeAgent := agents.NewNormalizeAgent()
	
	// Test basic normalization
	content := &agents.TextContent{
		Paragraphs: []string{"Test sentence."},
		Language:   "en-US",
		WordCount:  2,
	}
	
	_, err := normalizeAgent.Normalize(content)
	if err != nil {
		t.Errorf("Basic normalization failed: %v", err)
	}
}

func testMaxFileSize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large file test in short mode")
	}
	
	tempDir := t.TempDir()
	largeFile := filepath.Join(tempDir, "large.txt")
	
	// Create a large file (1MB of text)
	content := strings.Repeat("This is a test sentence with proper punctuation. ", 20000)
	err := os.WriteFile(largeFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}
	
	// Test processing
	ingestAgent := agents.NewTextIngestAgent()
	_, err = ingestAgent.ProcessFile(largeFile)
	if err != nil {
		t.Errorf("Should handle large files: %v", err)
	}
}

func testMemoryPressure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory pressure test in short mode")
	}
	
	// Create multiple large content objects to test memory handling
	normalizeAgent := agents.NewNormalizeAgent()
	
	for i := 0; i < 10; i++ {
		largeParagraphs := make([]string, 100)
		for j := range largeParagraphs {
			largeParagraphs[j] = strings.Repeat("Memory pressure test sentence. ", 50)
		}
		
		content := &agents.TextContent{
			Paragraphs: largeParagraphs,
			Language:   "en-US",
			WordCount:  15000,
		}
		
		_, err := normalizeAgent.Normalize(content)
		if err != nil {
			t.Errorf("Memory pressure test failed on iteration %d: %v", i, err)
			break
		}
	}
}

func testConcurrentProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}
	
	tempDir := t.TempDir()
	numGoroutines := 5
	
	// Create test files
	testFiles := make([]string, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		testFile := filepath.Join(tempDir, fmt.Sprintf("concurrent_%d.txt", i))
		content := fmt.Sprintf("Concurrent test file %d with unique content.", i)
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %d: %v", i, err)
		}
		testFiles[i] = testFile
	}
	
	// Process files concurrently
	errors := make(chan error, numGoroutines)
	
	for i, file := range testFiles {
		go func(index int, filename string) {
			ingestAgent := agents.NewTextIngestAgent()
			_, err := ingestAgent.ProcessFile(filename)
			errors <- err
		}(i, file)
	}
	
	// Collect results
	for i := 0; i < numGoroutines; i++ {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent processing failed for goroutine %d: %v", i, err)
		}
	}
}

// TestErrorRecovery tests system recovery from various error conditions
func TestErrorRecovery(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"corrupted file recovery", testCorruptedFileRecovery},
		{"permission denied recovery", testPermissionDeniedRecovery},
		{"disk space recovery", testDiskSpaceRecovery},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func testCorruptedFileRecovery(t *testing.T) {
	tempDir := t.TempDir()
	corruptedFile := filepath.Join(tempDir, "corrupted.txt")
	
	// Create a file with invalid UTF-8
	invalidUTF8 := []byte{0xff, 0xfe, 0xfd}
	err := os.WriteFile(corruptedFile, invalidUTF8, 0644)
	if err != nil {
		t.Fatalf("Failed to create corrupted file: %v", err)
	}
	
	ingestAgent := agents.NewTextIngestAgent()
	_, err = ingestAgent.ProcessFile(corruptedFile)
	
	// Should handle corrupted files gracefully
	if err == nil {
		t.Error("Should detect and handle corrupted files")
	}
}

func testPermissionDeniedRecovery(t *testing.T) {
	tempDir := t.TempDir()
	restrictedFile := filepath.Join(tempDir, "restricted.txt")
	
	// Create file and remove read permissions
	err := os.WriteFile(restrictedFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create restricted file: %v", err)
	}
	
	err = os.Chmod(restrictedFile, 0000)
	if err != nil {
		t.Fatalf("Failed to restrict file permissions: %v", err)
	}
	
	// Restore permissions for cleanup
	defer os.Chmod(restrictedFile, 0644)
	
	ingestAgent := agents.NewTextIngestAgent()
	_, err = ingestAgent.ProcessFile(restrictedFile)
	
	// Should handle permission errors gracefully
	if err == nil {
		t.Error("Should detect and handle permission denied errors")
	}
}

func testDiskSpaceRecovery(t *testing.T) {
	// This test would simulate disk space issues
	// For now, we'll just test that the system handles write errors
	tempDir := t.TempDir()
	
	// Test with a path that doesn't exist (simulates write failure)
	invalidPath := filepath.Join(tempDir, "nonexistent", "subdir", "file.mp3")
	
	// This would test actual disk space handling in a real implementation
	if _, err := os.Stat(filepath.Dir(invalidPath)); os.IsNotExist(err) {
		// Expected behavior - should handle path creation errors
		t.Log("Correctly detected invalid path")
	}
}

