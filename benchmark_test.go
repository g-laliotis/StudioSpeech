package main

import (
	"os"
	"path/filepath"
	"testing"
	"studiospeech/internal/agents"
)

// BenchmarkTextIngest benchmarks text file processing
func BenchmarkTextIngest(b *testing.B) {
	agent := agents.NewTextIngestAgent()
	
	// Create test file
	testFile := filepath.Join(b.TempDir(), "benchmark.txt")
	content := "This is a benchmark test file with multiple sentences. " +
		"It contains various punctuation marks! Does it work well? " +
		"We need to test performance with realistic content."
	
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.ProcessFile(testFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkNormalization benchmarks text normalization
func BenchmarkNormalization(b *testing.B) {
	agent := agents.NewNormalizeAgent()
	
	content := &agents.TextContent{
		Paragraphs: []string{
			"Dr. Smith said that 5 people attended the meeting at 3 PM.",
			"The results were amazing! We achieved 100% success rate.",
			"Please contact us at info@example.com for more details.",
		},
		Language:  "en-US",
		WordCount: 25,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.Normalize(content)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkVoiceSelection benchmarks voice catalog operations
func BenchmarkVoiceSelection(b *testing.B) {
	agent := agents.NewVoiceCatalogAgent("voices/catalog.json")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.SelectVoice("en-US", "female")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCacheOperations benchmarks cache hit/miss scenarios
func BenchmarkCacheOperations(b *testing.B) {
	cacheDir := b.TempDir()
	agent := agents.NewCacheAgent(cacheDir)
	
	testKey := "test-content-hash"
	testPath := filepath.Join(cacheDir, "test-audio.mp3")
	
	// Create test audio file
	if err := os.WriteFile(testPath, []byte("fake audio data"), 0644); err != nil {
		b.Fatal(err)
	}
	
	// Store in cache first
	agent.Store(testKey, testPath)
	
	b.Run("CacheHit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, found := agent.Get(testKey)
			if !found {
				b.Fatal("Expected cache hit")
			}
		}
	})
	
	b.Run("CacheMiss", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, found := agent.Get("non-existent-key")
			if found {
				b.Fatal("Expected cache miss")
			}
		}
	})
}

// BenchmarkPipelineEnd2End benchmarks the complete pipeline
func BenchmarkPipelineEnd2End(b *testing.B) {
	// Skip if system dependencies not available
	if !isSystemReady() {
		b.Skip("System dependencies not available")
	}
	
	testFile := filepath.Join(b.TempDir(), "pipeline-test.txt")
	outputFile := filepath.Join(b.TempDir(), "output.mp3")
	
	content := "Hello world. This is a comprehensive benchmark test. " +
		"It measures the performance of the entire TTS pipeline."
	
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Run pipeline (simplified version for benchmarking)
		if err := runSimplifiedPipeline(testFile, outputFile); err != nil {
			b.Fatal(err)
		}
		
		// Clean up output for next iteration
		os.Remove(outputFile)
	}
}

// BenchmarkMemoryUsage benchmarks memory allocation patterns
func BenchmarkMemoryUsage(b *testing.B) {
	agent := agents.NewNormalizeAgent()
	
	// Large text content to test memory usage
	largeContent := &agents.TextContent{
		Paragraphs: make([]string, 100),
		Language:   "en-US",
		WordCount:  1000,
	}
	
	for i := range largeContent.Paragraphs {
		largeContent.Paragraphs[i] = "This is paragraph number " + 
			"with multiple sentences and various punctuation marks! " +
			"Does it handle memory efficiently? We need to test this carefully."
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, err := agent.Normalize(largeContent)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Helper functions for benchmarks

func isSystemReady() bool {
	// Check if basic system requirements are met
	envAgent := agents.NewEnvironmentAgent()
	return envAgent.Validate() == nil
}

func runSimplifiedPipeline(inputFile, outputFile string) error {
	// Simplified pipeline for benchmarking
	// This would run the actual pipeline in a real implementation
	
	// Text ingestion
	ingestAgent := agents.NewTextIngestAgent()
	content, err := ingestAgent.ProcessFile(inputFile)
	if err != nil {
		return err
	}
	
	// Normalization
	normalizeAgent := agents.NewNormalizeAgent()
	normalized, err := normalizeAgent.Normalize(content)
	if err != nil {
		return err
	}
	
	// For benchmarking, we'll just create a dummy output file
	return os.WriteFile(outputFile, []byte("dummy audio data"), 0644)
}