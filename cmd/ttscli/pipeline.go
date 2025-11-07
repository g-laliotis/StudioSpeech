package main

import (
	"fmt"
	"os"
	"path/filepath"

	"studiospeech/internal/agents"
)

// executeSynthesisPipeline runs the complete TTS pipeline
func executeSynthesisPipeline() error {
	fmt.Println("🔄 Starting synthesis pipeline...")

	// Step 1: Text Ingestion
	fmt.Printf("📖 Reading input file: %s\n", inputFile)
	textAgent := agents.NewTextIngestAgent()
	content, err := textAgent.ProcessFile(inputFile)
	if err != nil {
		return fmt.Errorf("text ingestion failed: %w", err)
	}

	if err := textAgent.ValidateContent(content); err != nil {
		return fmt.Errorf("content validation failed: %w", err)
	}

	fmt.Printf("   ✓ Processed %d paragraphs, %d words\n", len(content.Paragraphs), content.WordCount)
	fmt.Printf("   ✓ Detected language: %s\n", content.Language)

	// Override language if specified
	if language != "auto" {
		content.Language = language
		fmt.Printf("   ✓ Language override: %s\n", language)
	}

	// Step 2: Voice Selection
	fmt.Printf("🎭 Selecting voice...\n")
	catalogPath := filepath.Join("voices", "catalog.json")
	voiceAgent := agents.NewVoiceCatalogAgent(catalogPath)

	if err := voiceAgent.LoadCatalog(); err != nil {
		return fmt.Errorf("voice catalog loading failed: %w", err)
	}

	selectedVoice, err := voiceAgent.SelectVoice(content.Language, voiceID, gender)
	if err != nil {
		return fmt.Errorf("voice selection failed: %w", err)
	}

	fmt.Printf("   ✓ Selected voice: %s (%s %s)\n", selectedVoice.ID, selectedVoice.Gender, selectedVoice.Style)

	// Step 3: Text Normalization
	fmt.Printf("🔧 Normalizing text...\n")
	normalizeAgent := agents.NewNormalizeAgent()
	normalized, err := normalizeAgent.Normalize(content)
	if err != nil {
		return fmt.Errorf("text normalization failed: %w", err)
	}

	if err := normalizeAgent.ValidateNormalizedText(normalized); err != nil {
		return fmt.Errorf("normalized text validation failed: %w", err)
	}

	fmt.Printf("   ✓ Generated %d sentences\n", len(normalized.Sentences))

	// Step 4: Synthesis (Dry Run for now)
	fmt.Printf("🎤 Synthesizing speech...\n")

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "studiospeech_*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize synthesis agent
	synthAgent := agents.NewSynthAgent("piper", tempDir)
	synthAgent.SetDryRun(false) // Use real synthesis with fallback

	// Set synthesis parameters
	params := &agents.SynthParams{
		Speed:   speed,
		Noise:   noise,
		NoiseW:  noisew,
		Speaker: 0,
	}

	result, err := synthAgent.Synthesize(normalized, selectedVoice, params)
	if err != nil {
		return fmt.Errorf("synthesis failed: %w", err)
	}

	fmt.Printf("   ✓ Generated audio: %s\n", result.OutputPath)
	fmt.Printf("   ✓ Sample rate: %d Hz, Channels: %d\n", result.SampleRate, result.Channels)

	// Step 5: Check Cache
	fmt.Printf("💾 Checking cache...\n")
	cacheDir := filepath.Join(os.TempDir(), "studiospeech_cache")
	cacheAgent := agents.NewCacheAgent(cacheDir)
	if err := cacheAgent.Initialize(); err != nil {
		return fmt.Errorf("cache initialization failed: %w", err)
	}

	// Generate cache key
	postParams := &agents.PostProcessParams{
		Format:       agents.AudioFormat(format),
		SampleRate:   sampleRate,
		Bitrate:      bitrate,
		LoudnessLUFS: -16.0,
	}

	cacheKey := cacheAgent.GenerateKey(content, selectedVoice, params, postParams)

	// Check for cache hit
	if entry, err := cacheAgent.Get(cacheKey); err == nil && entry != nil {
		fmt.Printf("   ✅ Cache hit! Using cached audio\n")
		fmt.Printf("   ✓ Cached file: %s\n", entry.FilePath)

		// Copy cached file to output location
		if err := copyFile(entry.FilePath, outputFile); err != nil {
			return fmt.Errorf("failed to copy cached file: %w", err)
		}
		return nil
	}

	fmt.Printf("   ⚠️  Cache miss - will synthesize and cache result\n")

	// Step 6: Post-processing
	fmt.Printf("🎵 Post-processing...\n")
	postAgent := agents.NewPostProcessAgent("ffmpeg", tempDir)
	postAgent.SetDryRun(false) // Use real post-processing

	postResult, err := postAgent.Process(result.OutputPath, outputFile, postParams)
	if err != nil {
		return fmt.Errorf("post-processing failed: %w", err)
	}

	fmt.Printf("   ✓ Processed audio: %s\n", postResult.OutputPath)
	fmt.Printf("   ✓ Format: %s, Sample rate: %d Hz\n", postResult.Format, postResult.SampleRate)

	// Step 7: Cache result
	fmt.Printf("💾 Caching result...\n")
	metadata := map[string]interface{}{
		"voice":      selectedVoice.ID,
		"language":   content.Language,
		"format":     string(postParams.Format),
		"word_count": content.WordCount,
	}

	if err := cacheAgent.Put(cacheKey, postResult.OutputPath, metadata); err != nil {
		fmt.Printf("   ⚠️  Failed to cache result: %v\n", err)
	} else {
		fmt.Printf("   ✓ Result cached for future use\n")
	}

	fmt.Printf("   ⚠️  This is a dry-run - no actual audio generated\n")
	fmt.Printf("   ⚠️  Install Piper TTS and FFmpeg to enable real synthesis\n")

	return nil
}
