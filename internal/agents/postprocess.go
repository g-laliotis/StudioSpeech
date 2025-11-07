package agents

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// AudioFormat represents output audio format
type AudioFormat string

const (
	FormatWAV AudioFormat = "wav"
	FormatMP3 AudioFormat = "mp3"
)

// PostProcessParams contains audio processing parameters
type PostProcessParams struct {
	Format     AudioFormat
	SampleRate int     // Target sample rate (Hz)
	Bitrate    int     // MP3 bitrate (kbps)
	LoudnessLUFS float64 // Target loudness (-16 to -14 LUFS)
}

// PostProcessResult contains processing output information
type PostProcessResult struct {
	OutputPath   string
	Format       AudioFormat
	SampleRate   int
	Channels     int
	Duration     float64
	FileSize     int64
}

// PostProcessAgent handles audio post-processing using FFmpeg
type PostProcessAgent struct {
	ffmpegPath string
	tempDir    string
	dryRun     bool
}

// NewPostProcessAgent creates a new post-processing agent
func NewPostProcessAgent(ffmpegPath, tempDir string) *PostProcessAgent {
	return &PostProcessAgent{
		ffmpegPath: ffmpegPath,
		tempDir:    tempDir,
		dryRun:     false,
	}
}

// SetDryRun enables/disables dry-run mode for testing
func (p *PostProcessAgent) SetDryRun(enabled bool) {
	p.dryRun = enabled
}

// Process converts and normalizes audio using FFmpeg
func (p *PostProcessAgent) Process(inputPath, outputPath string, params *PostProcessParams) (*PostProcessResult, error) {
	if params == nil {
		params = p.getDefaultParams()
	}
	
	if err := p.validateParams(params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}
	
	// Check input file exists (skip in dry-run mode)
	if !p.dryRun {
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("input file not found: %s", inputPath)
		}
	}
	
	// Build FFmpeg command
	cmd := p.buildFFmpegCommand(inputPath, outputPath, params)
	
	if p.dryRun {
		return &PostProcessResult{
			OutputPath: outputPath,
			Format:     params.Format,
			SampleRate: params.SampleRate,
			Channels:   1,
			Duration:   0,
			FileSize:   0,
		}, nil
	}
	
	// Execute FFmpeg
	if err := p.executeFFmpeg(cmd); err != nil {
		return nil, fmt.Errorf("ffmpeg processing failed: %w", err)
	}
	
	// Get output file info
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get output file info: %w", err)
	}
	
	return &PostProcessResult{
		OutputPath: outputPath,
		Format:     params.Format,
		SampleRate: params.SampleRate,
		Channels:   1,
		Duration:   0, // Would need ffprobe to get actual duration
		FileSize:   fileInfo.Size(),
	}, nil
}

// buildFFmpegCommand constructs the FFmpeg command line
func (p *PostProcessAgent) buildFFmpegCommand(inputPath, outputPath string, params *PostProcessParams) *exec.Cmd {
	args := []string{
		"-i", inputPath,
		"-y", // Overwrite output file
	}
	
	// Audio processing filters
	var filters []string
	
	// Resample to target sample rate and convert to mono
	filters = append(filters, fmt.Sprintf("aresample=%d", params.SampleRate))
	filters = append(filters, "pan=mono|c0=0.5*c0+0.5*c1")
	
	// Loudness normalization (EBU R128)
	if params.LoudnessLUFS != 0 {
		filters = append(filters, fmt.Sprintf("loudnorm=I=%.1f:TP=-1.0:LRA=7.0", params.LoudnessLUFS))
	}
	
	// Apply filters
	if len(filters) > 0 {
		args = append(args, "-af", strings.Join(filters, ","))
	}
	
	// Format-specific options
	switch params.Format {
	case FormatMP3:
		args = append(args, "-codec:a", "libmp3lame")
		args = append(args, "-b:a", fmt.Sprintf("%dk", params.Bitrate))
		args = append(args, "-ar", strconv.Itoa(params.SampleRate))
	case FormatWAV:
		args = append(args, "-codec:a", "pcm_s16le")
		args = append(args, "-ar", strconv.Itoa(params.SampleRate))
	}
	
	args = append(args, outputPath)
	
	return exec.Command(p.ffmpegPath, args...)
}

// executeFFmpeg runs FFmpeg command
func (p *PostProcessAgent) executeFFmpeg(cmd *exec.Cmd) error {
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// getDefaultParams returns safe default processing parameters
func (p *PostProcessAgent) getDefaultParams() *PostProcessParams {
	return &PostProcessParams{
		Format:       FormatMP3,
		SampleRate:   48000,
		Bitrate:      192,
		LoudnessLUFS: -16.0, // YouTube-friendly loudness
	}
}

// validateParams checks if processing parameters are valid
func (p *PostProcessAgent) validateParams(params *PostProcessParams) error {
	if params.Format != FormatWAV && params.Format != FormatMP3 {
		return fmt.Errorf("unsupported format: %s", params.Format)
	}
	
	if params.SampleRate < 8000 || params.SampleRate > 192000 {
		return fmt.Errorf("sample rate must be between 8000 and 192000 Hz, got %d", params.SampleRate)
	}
	
	if params.Format == FormatMP3 {
		if params.Bitrate < 64 || params.Bitrate > 320 {
			return fmt.Errorf("MP3 bitrate must be between 64 and 320 kbps, got %d", params.Bitrate)
		}
	}
	
	if params.LoudnessLUFS != 0 && (params.LoudnessLUFS < -30 || params.LoudnessLUFS > -6) {
		return fmt.Errorf("loudness must be between -30 and -6 LUFS, got %.1f", params.LoudnessLUFS)
	}
	
	return nil
}

// GetCommandLine returns the command line that would be executed (for testing)
func (p *PostProcessAgent) GetCommandLine(inputPath, outputPath string, params *PostProcessParams) string {
	if params == nil {
		params = p.getDefaultParams()
	}
	
	cmd := p.buildFFmpegCommand(inputPath, outputPath, params)
	return strings.Join(append([]string{cmd.Path}, cmd.Args[1:]...), " ")
}

// CleanupTempFiles removes temporary processing files
func (p *PostProcessAgent) CleanupTempFiles(result *PostProcessResult) error {
	if result != nil && result.OutputPath != "" {
		if err := os.Remove(result.OutputPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to cleanup temp file %s: %w", result.OutputPath, err)
		}
	}
	return nil
}