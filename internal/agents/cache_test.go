package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCacheAgent_GenerateKey(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	agent := NewCacheAgent(tempDir)
	
	content := &TextContent{
		Paragraphs: []string{"Hello world", "Test content"},
		Language:   "en-US",
	}
	
	voice := &Voice{
		ID: "test_voice",
	}
	
	synthParams := &SynthParams{
		Speed:  1.0,
		Noise:  0.5,
		NoiseW: 0.8,
	}
	
	postParams := &PostProcessParams{
		Format:     FormatMP3,
		SampleRate: 48000,
		Bitrate:    192,
	}
	
	key1 := agent.GenerateKey(content, voice, synthParams, postParams)
	key2 := agent.GenerateKey(content, voice, synthParams, postParams)
	
	// Same inputs should generate same key
	if key1 != key2 {
		t.Error("Same inputs should generate same cache key")
	}
	
	// Different content should generate different key
	content2 := &TextContent{
		Paragraphs: []string{"Different content"},
		Language:   "en-US",
	}
	
	key3 := agent.GenerateKey(content2, voice, synthParams, postParams)
	if key1 == key3 {
		t.Error("Different content should generate different cache key")
	}
}

func TestCacheAgent_PutAndGet(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	agent := NewCacheAgent(tempDir)
	if err := agent.Initialize(); err != nil {
		t.Fatalf("Failed to initialize cache: %v", err)
	}
	
	// Create a test file
	testFile := filepath.Join(tempDir, "test.mp3")
	if err := os.WriteFile(testFile, []byte("test audio data"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	key := "test_key_123"
	metadata := map[string]interface{}{
		"voice": "test_voice",
		"format": "mp3",
	}
	
	// Put file in cache
	if err := agent.Put(key, testFile, metadata); err != nil {
		t.Fatalf("Failed to put file in cache: %v", err)
	}
	
	// Get file from cache
	entry, err := agent.Get(key)
	if err != nil {
		t.Fatalf("Failed to get file from cache: %v", err)
	}
	
	if entry == nil {
		t.Fatal("Cache entry should not be nil")
	}
	
	if entry.Key != key {
		t.Errorf("Expected key %s, got %s", key, entry.Key)
	}
	
	// Check cached file exists
	if _, err := os.Stat(entry.FilePath); os.IsNotExist(err) {
		t.Error("Cached file should exist")
	}
}

func TestCacheAgent_CacheMiss(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	agent := NewCacheAgent(tempDir)
	if err := agent.Initialize(); err != nil {
		t.Fatalf("Failed to initialize cache: %v", err)
	}
	
	// Try to get non-existent key
	entry, err := agent.Get("non_existent_key")
	if err != nil {
		t.Fatalf("Get should not return error for cache miss: %v", err)
	}
	
	if entry != nil {
		t.Error("Cache miss should return nil entry")
	}
}

func TestCacheAgent_Stats(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	agent := NewCacheAgent(tempDir)
	if err := agent.Initialize(); err != nil {
		t.Fatalf("Failed to initialize cache: %v", err)
	}
	
	stats := agent.Stats()
	
	if entries, ok := stats["entries"].(int); !ok || entries != 0 {
		t.Errorf("Expected 0 entries, got %v", stats["entries"])
	}
	
	if totalSize, ok := stats["total_size"].(int64); !ok || totalSize != 0 {
		t.Errorf("Expected 0 total size, got %v", stats["total_size"])
	}
}