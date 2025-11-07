package agents

import (
	"fmt"
	"os/exec"
	"runtime"
)

// MacOSTTSAgent provides fallback TTS using macOS built-in voices
type MacOSTTSAgent struct {
	tempDir string
}

// NewMacOSTTSAgent creates a new macOS TTS agent
func NewMacOSTTSAgent(tempDir string) *MacOSTTSAgent {
	return &MacOSTTSAgent{
		tempDir: tempDir,
	}
}

// IsAvailable checks if macOS TTS is available
func (m *MacOSTTSAgent) IsAvailable() bool {
	if runtime.GOOS != "darwin" {
		return false
	}

	// Check if 'say' command exists
	_, err := exec.LookPath("say")
	return err == nil
}

// Synthesize converts text to speech using macOS say command
func (m *MacOSTTSAgent) Synthesize(text, outputPath, gender, language string) error {
	if !m.IsAvailable() {
		return fmt.Errorf("macOS TTS not available")
	}

	// Select voice based on gender and language
	voice := m.selectVoice(gender, language)

	// Generate to AIFF first (macOS native format) with natural speech rate
	aiffPath := outputPath + ".aiff"
	rate := m.calculateSpeechRate(language)
	cmd := exec.Command("say", "-v", voice, "-r", fmt.Sprintf("%.0f", rate), "-o", aiffPath, text)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("macOS TTS failed: %w\nOutput: %s", err, string(output))
	}

	// Convert AIFF to WAV using ffmpeg (if available)
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		convertCmd := exec.Command("ffmpeg", "-i", aiffPath, "-y", outputPath)
		if convertErr := convertCmd.Run(); convertErr == nil {
			// Remove temporary AIFF file
			exec.Command("rm", aiffPath).Run()
			return nil
		}
	}

	// If conversion fails, just rename AIFF to output path
	exec.Command("mv", aiffPath, outputPath).Run()
	return nil
}

// selectVoice chooses appropriate macOS voice based on gender and language
func (m *MacOSTTSAgent) selectVoice(gender, language string) string {
	// Check for Greek language
	if language == "el-GR" {
		return "Melina" // Greek female voice
	}

	// Default to English voices
	switch gender {
	case "male":
		return "Alex" // Default male voice
	case "female":
		return "Samantha" // Default female voice
	default:
		return "Samantha" // Default to female
	}
}

// calculateSpeechRate returns optimal speech rate for natural delivery
func (m *MacOSTTSAgent) calculateSpeechRate(language string) float64 {
	// Slower rate for better pronunciation and natural pauses
	switch language {
	case "el-GR":
		return 160 // Slightly slower for Greek
	default:
		return 175 // Natural English rate
	}
}

// GetAvailableVoices returns list of available macOS voices
func (m *MacOSTTSAgent) GetAvailableVoices() ([]string, error) {
	cmd := exec.Command("say", "-v", "?")
	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get voices: %w", err)
	}

	// Parse voice list (simplified)
	return []string{"Alex (male)", "Samantha (female)", "Victoria (female)"}, nil
}
