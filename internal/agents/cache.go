package agents

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// CacheEntry represents a cached synthesis result
type CacheEntry struct {
	Key       string    `json:"key"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	FileSize  int64     `json:"file_size"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// CacheIndex maintains the cache index
type CacheIndex struct {
	Entries map[string]*CacheEntry `json:"entries"`
	Version string                 `json:"version"`
}

// CacheAgent handles synthesis result caching
type CacheAgent struct {
	cacheDir   string
	indexPath  string
	index      *CacheIndex
	maxAge     time.Duration
	maxSize    int64
}

// NewCacheAgent creates a new cache agent
func NewCacheAgent(cacheDir string) *CacheAgent {
	return &CacheAgent{
		cacheDir:  cacheDir,
		indexPath: filepath.Join(cacheDir, "index.json"),
		maxAge:    24 * time.Hour, // 24 hours default
		maxSize:   1024 * 1024 * 1024, // 1GB default
	}
}

// Initialize creates cache directory and loads index
func (c *CacheAgent) Initialize() error {
	// Create cache directory
	if err := os.MkdirAll(c.cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}
	
	// Load or create index
	if err := c.loadIndex(); err != nil {
		return fmt.Errorf("failed to load cache index: %w", err)
	}
	
	return nil
}

// GenerateKey creates a cache key from content and parameters
func (c *CacheAgent) GenerateKey(content *TextContent, voice *Voice, synthParams *SynthParams, postParams *PostProcessParams) string {
	hasher := sha256.New()
	
	// Hash text content
	for _, paragraph := range content.Paragraphs {
		hasher.Write([]byte(paragraph))
	}
	
	// Hash voice ID
	hasher.Write([]byte(voice.ID))
	
	// Hash synthesis parameters
	if synthParams != nil {
		hasher.Write([]byte(fmt.Sprintf("%.3f-%.3f-%.3f-%d", 
			synthParams.Speed, synthParams.Noise, synthParams.NoiseW, synthParams.Speaker)))
	}
	
	// Hash post-processing parameters
	if postParams != nil {
		hasher.Write([]byte(fmt.Sprintf("%s-%d-%d-%.1f", 
			postParams.Format, postParams.SampleRate, postParams.Bitrate, postParams.LoudnessLUFS)))
	}
	
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// Get retrieves cached result if available
func (c *CacheAgent) Get(key string) (*CacheEntry, error) {
	if c.index == nil {
		return nil, fmt.Errorf("cache not initialized")
	}
	
	entry, exists := c.index.Entries[key]
	if !exists {
		return nil, nil // Cache miss
	}
	
	// Check if file still exists
	if _, err := os.Stat(entry.FilePath); os.IsNotExist(err) {
		// File missing, remove from index
		delete(c.index.Entries, key)
		c.saveIndex()
		return nil, nil
	}
	
	// Check if entry is too old
	if time.Since(entry.CreatedAt) > c.maxAge {
		c.Remove(key)
		return nil, nil
	}
	
	return entry, nil
}

// Put stores a result in cache
func (c *CacheAgent) Put(key, filePath string, metadata map[string]interface{}) error {
	if c.index == nil {
		return fmt.Errorf("cache not initialized")
	}
	
	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	
	// Create cache file path
	cacheFilePath := filepath.Join(c.cacheDir, key+filepath.Ext(filePath))
	
	// Copy file to cache
	if err := c.copyFile(filePath, cacheFilePath); err != nil {
		return fmt.Errorf("failed to copy file to cache: %w", err)
	}
	
	// Create cache entry
	entry := &CacheEntry{
		Key:       key,
		FilePath:  cacheFilePath,
		CreatedAt: time.Now(),
		FileSize:  fileInfo.Size(),
		Metadata:  metadata,
	}
	
	// Add to index
	c.index.Entries[key] = entry
	
	// Save index
	if err := c.saveIndex(); err != nil {
		return fmt.Errorf("failed to save cache index: %w", err)
	}
	
	return nil
}

// Remove deletes a cache entry
func (c *CacheAgent) Remove(key string) error {
	if c.index == nil {
		return fmt.Errorf("cache not initialized")
	}
	
	entry, exists := c.index.Entries[key]
	if !exists {
		return nil // Already removed
	}
	
	// Remove file
	if err := os.Remove(entry.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cache file: %w", err)
	}
	
	// Remove from index
	delete(c.index.Entries, key)
	
	// Save index
	return c.saveIndex()
}

// Prune removes old or large cache entries
func (c *CacheAgent) Prune() error {
	if c.index == nil {
		return fmt.Errorf("cache not initialized")
	}
	
	var totalSize int64
	var toRemove []string
	
	// Calculate total size and find old entries
	for key, entry := range c.index.Entries {
		totalSize += entry.FileSize
		
		// Mark old entries for removal
		if time.Since(entry.CreatedAt) > c.maxAge {
			toRemove = append(toRemove, key)
		}
	}
	
	// Remove old entries
	for _, key := range toRemove {
		c.Remove(key)
		totalSize -= c.index.Entries[key].FileSize
	}
	
	// If still over size limit, remove oldest entries
	if totalSize > c.maxSize {
		// Sort by creation time and remove oldest
		// Simplified: just remove entries until under limit
		for key, entry := range c.index.Entries {
			if totalSize <= c.maxSize {
				break
			}
			c.Remove(key)
			totalSize -= entry.FileSize
		}
	}
	
	return nil
}

// Stats returns cache statistics
func (c *CacheAgent) Stats() map[string]interface{} {
	if c.index == nil {
		return map[string]interface{}{"error": "cache not initialized"}
	}
	
	var totalSize int64
	entryCount := len(c.index.Entries)
	
	for _, entry := range c.index.Entries {
		totalSize += entry.FileSize
	}
	
	return map[string]interface{}{
		"entries":    entryCount,
		"total_size": totalSize,
		"cache_dir":  c.cacheDir,
	}
}

// loadIndex loads the cache index from disk
func (c *CacheAgent) loadIndex() error {
	// Initialize empty index if file doesn't exist
	if _, err := os.Stat(c.indexPath); os.IsNotExist(err) {
		c.index = &CacheIndex{
			Entries: make(map[string]*CacheEntry),
			Version: "1.0",
		}
		return c.saveIndex()
	}
	
	// Load existing index
	file, err := os.Open(c.indexPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	c.index = &CacheIndex{}
	if err := json.NewDecoder(file).Decode(c.index); err != nil {
		// If index is corrupted, start fresh
		c.index = &CacheIndex{
			Entries: make(map[string]*CacheEntry),
			Version: "1.0",
		}
		return c.saveIndex()
	}
	
	// Ensure entries map is initialized
	if c.index.Entries == nil {
		c.index.Entries = make(map[string]*CacheEntry)
	}
	
	return nil
}

// saveIndex saves the cache index to disk
func (c *CacheAgent) saveIndex() error {
	file, err := os.Create(c.indexPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c.index)
}

// copyFile copies a file from src to dst
func (c *CacheAgent) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	_, err = io.Copy(dstFile, srcFile)
	return err
}