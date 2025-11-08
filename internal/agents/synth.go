package agents

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// SynthParams contains synthesis parameters
type SynthParams struct {
	Speed   float64 // Speech speed multiplier (0.5-2.0)
	Noise   float64 // Noise level for naturalness (0.0-1.0)
	NoiseW  float64 // Noise width for variation (0.0-1.0)
	Speaker int     // Speaker ID for multi-speaker models
}

// SynthResult contains synthesis output information
type SynthResult struct {
	OutputPath string
	Duration   time.Duration
	SampleRate int
	Channels   int
	FileSize   int64
}

// SynthAgent handles text-to-speech synthesis using Piper
type SynthAgent struct {
	piperPath string
	tempDir   string
	dryRun    bool
}

// NewSynthAgent creates a new synthesis agent
func NewSynthAgent(piperPath, tempDir string) *SynthAgent {
	return &SynthAgent{
		piperPath: piperPath,
		tempDir:   tempDir,
		dryRun:    false,
	}
}

// SetDryRun enables/disables dry-run mode for testing
func (s *SynthAgent) SetDryRun(enabled bool) {
	s.dryRun = enabled
}

// Synthesize converts normalized text to speech using Piper
func (s *SynthAgent) Synthesize(normalized *NormalizedText, voice *Voice, params *SynthParams) (*SynthResult, error) {
	if normalized == nil {
		return nil, fmt.Errorf("normalized text is nil")
	}

	if voice == nil {
		return nil, fmt.Errorf("voice is nil")
	}

	if params == nil {
		params = s.getDefaultParams()
	}

	// Validate parameters
	if err := s.validateParams(params); err != nil {
		return nil, fmt.Errorf("invalid synthesis parameters: %w", err)
	}

	// Check if voice model file exists (skip for macOS voices and dry-run mode)
	if !s.dryRun && !s.isMacOSVoice(voice) {
		if _, err := os.Stat(voice.Path); os.IsNotExist(err) {
			return nil, fmt.Errorf("voice model file not found: %s", voice.Path)
		}
	}

	// Create temporary output file
	outputPath := filepath.Join(s.tempDir, fmt.Sprintf("synth_%d.wav", time.Now().UnixNano()))

	// Combine sentences with proper pauses between them
	text := strings.Join(normalized.Sentences, ". ")

	// Build Piper command
	cmd := s.buildPiperCommand(voice.Path, outputPath, params)

	if s.dryRun {
		// Return command for testing without execution
		return &SynthResult{
			OutputPath: outputPath,
			Duration:   0,
			SampleRate: voice.SampleRate,
			Channels:   1,
			FileSize:   0,
		}, nil
	}

	// Execute synthesis
	startTime := time.Now()

	// Use macOS TTS for macOS voices, Piper for others
	if s.isMacOSVoice(voice) {
		macTTS := NewMacOSTTSAgent(s.tempDir)
		if macTTS.IsAvailable() {
			if err := macTTS.Synthesize(text, outputPath, voice.Gender, normalized.Language); err != nil {
				return nil, fmt.Errorf("macOS TTS synthesis failed: %w", err)
			}
		} else {
			return nil, fmt.Errorf("macOS TTS not available")
		}
	} else {
		// Try Piper first, fallback to macOS TTS if Piper fails
		err := s.executePiper(cmd, text)
		if err != nil {
			// Try macOS TTS fallback
			macTTS := NewMacOSTTSAgent(s.tempDir)
			if macTTS.IsAvailable() {
				gender := "female"
				if params.Speaker > 0 {
					gender = "male"
				}
				if err := macTTS.Synthesize(text, outputPath, gender, normalized.Language); err != nil {
					return nil, fmt.Errorf("both Piper and macOS TTS failed: piper=%v, macos=%v", err, err)
				}
			} else {
				return nil, fmt.Errorf("piper synthesis failed and no fallback available: %w", err)
			}
		}
	}

	duration := time.Since(startTime)

	// Get output file info
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get output file info: %w", err)
	}

	return &SynthResult{
		OutputPath: outputPath,
		Duration:   duration,
		SampleRate: voice.SampleRate,
		Channels:   1,
		FileSize:   fileInfo.Size(),
	}, nil
}

// buildPiperCommand constructs the Piper command line
func (s *SynthAgent) buildPiperCommand(modelPath, outputPath string, params *SynthParams) *exec.Cmd {
	args := []string{
		"--model", modelPath,
		"--output_file", outputPath,
	}

	// Convert speed to length_scale (inverse relationship)
	lengthScale := 1.0 / params.Speed
	args = append(args, "--length_scale", fmt.Sprintf("%.3f", lengthScale))

	// Add noise parameters
	args = append(args, "--noise_scale", fmt.Sprintf("%.3f", params.Noise))
	args = append(args, "--noise_w", fmt.Sprintf("%.3f", params.NoiseW))

	// Add speaker if specified
	if params.Speaker > 0 {
		args = append(args, "--speaker", strconv.Itoa(params.Speaker))
	}

	return exec.Command(s.piperPath, args...)
}

// executePiper runs Piper with the given text input
func (s *SynthAgent) executePiper(cmd *exec.Cmd, text string) error {
	// Set up stdin pipe for text input
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start piper: %w", err)
	}

	// Write text to stdin
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, text)
	}()

	// Wait for completion
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("piper execution failed: %w", err)
	}

	return nil
}

// getDefaultParams returns safe default synthesis parameters
func (s *SynthAgent) getDefaultParams() *SynthParams {
	return &SynthParams{
		Speed:   1.03,  // Slightly faster than normal
		Noise:   0.667, // Natural variation
		NoiseW:  0.8,   // Moderate width
		Speaker: 0,     // Default speaker
	}
}

// validateParams checks if synthesis parameters are within valid ranges
func (s *SynthAgent) validateParams(params *SynthParams) error {
	if params.Speed < 0.5 || params.Speed > 2.0 {
		return fmt.Errorf("speed must be between 0.5 and 2.0, got %.2f", params.Speed)
	}

	if params.Noise < 0.0 || params.Noise > 1.0 {
		return fmt.Errorf("noise must be between 0.0 and 1.0, got %.3f", params.Noise)
	}

	if params.NoiseW < 0.0 || params.NoiseW > 1.0 {
		return fmt.Errorf("noiseW must be between 0.0 and 1.0, got %.3f", params.NoiseW)
	}

	if params.Speaker < 0 {
		return fmt.Errorf("speaker must be >= 0, got %d", params.Speaker)
	}

	return nil
}

// GetCommandLine returns the command line that would be executed (for testing)
func (s *SynthAgent) GetCommandLine(voice *Voice, params *SynthParams, outputPath string) string {
	if params == nil {
		params = s.getDefaultParams()
	}

	cmd := s.buildPiperCommand(voice.Path, outputPath, params)
	return strings.Join(append([]string{cmd.Path}, cmd.Args[1:]...), " ")
}

// isMacOSVoice checks if a voice is a macOS system voice
func (s *SynthAgent) isMacOSVoice(voice *Voice) bool {
	// macOS voices have simple names like "Alex", "Samantha", "Melina"
	// and don't have file extensions
	return !strings.Contains(voice.Path, "/") && !strings.Contains(voice.Path, ".")
}

// CleanupTempFiles removes temporary synthesis files
func (s *SynthAgent) CleanupTempFiles(result *SynthResult) error {
	if result != nil && result.OutputPath != "" {
		if err := os.Remove(result.OutputPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to cleanup temp file %s: %w", result.OutputPath, err)
		}
	}
	return nil
}
