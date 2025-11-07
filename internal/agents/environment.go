package agents

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// EnvironmentInfo contains system environment details
type EnvironmentInfo struct {
	PiperPath     string
	PiperVersion  string
	FFmpegPath    string
	FFmpegVersion string
	OS            string
	Arch          string
	GoVersion     string
	HasPiper      bool
	HasFFmpeg     bool
}

// EnvironmentAgent handles system requirements validation
type EnvironmentAgent struct{}

// NewEnvironmentAgent creates a new environment validation agent
func NewEnvironmentAgent() *EnvironmentAgent {
	return &EnvironmentAgent{}
}

// Check validates system requirements and returns environment info
func (e *EnvironmentAgent) Check() (*EnvironmentInfo, error) {
	info := &EnvironmentInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
	}

	// Check Piper TTS
	if err := e.checkPiper(info); err != nil {
		// Piper not found, but continue checking other components
	}

	// Check FFmpeg
	if err := e.checkFFmpeg(info); err != nil {
		// FFmpeg not found, but continue checking other components
	}

	return info, nil
}

// checkPiper validates Piper TTS installation
func (e *EnvironmentAgent) checkPiper(info *EnvironmentInfo) error {
	// Try to find piper in PATH
	path, err := exec.LookPath("piper")
	if err != nil {
		info.HasPiper = false
		return fmt.Errorf("piper not found in PATH")
	}

	info.PiperPath = path

	// Get version
	cmd := exec.Command("piper", "--version")
	output, err := cmd.Output()
	if err != nil {
		info.HasPiper = false
		return fmt.Errorf("failed to get piper version: %w", err)
	}

	info.PiperVersion = strings.TrimSpace(string(output))
	info.HasPiper = true

	return nil
}

// checkFFmpeg validates FFmpeg installation
func (e *EnvironmentAgent) checkFFmpeg(info *EnvironmentInfo) error {
	// Try to find ffmpeg in PATH
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		info.HasFFmpeg = false
		return fmt.Errorf("ffmpeg not found in PATH")
	}

	info.FFmpegPath = path

	// Get version
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.Output()
	if err != nil {
		info.HasFFmpeg = false
		return fmt.Errorf("failed to get ffmpeg version: %w", err)
	}

	// Extract version from first line (e.g., "ffmpeg version 4.4.2")
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		fields := strings.Fields(lines[0])
		if len(fields) >= 3 {
			info.FFmpegVersion = fields[2]
		} else {
			info.FFmpegVersion = "unknown"
		}
	}

	info.HasFFmpeg = true
	return nil
}

// GetInstallGuide returns OS-specific installation instructions
func (e *EnvironmentAgent) GetInstallGuide(missing []string) string {
	var guide strings.Builder

	guide.WriteString("📋 Installation Guide:\n\n")

	for _, tool := range missing {
		switch tool {
		case "piper":
			guide.WriteString("Piper TTS:\n")
			switch runtime.GOOS {
			case "darwin":
				guide.WriteString("  macOS: brew install piper-tts\n")
				guide.WriteString("  Or download: https://github.com/rhasspy/piper/releases\n")
			case "windows":
				guide.WriteString("  Windows: choco install piper-tts\n")
				guide.WriteString("  Or download: https://github.com/rhasspy/piper/releases\n")
			default:
				guide.WriteString("  Linux: Download from https://github.com/rhasspy/piper/releases\n")
				guide.WriteString("  Extract and add to PATH\n")
			}
			guide.WriteString("\n")

		case "ffmpeg":
			guide.WriteString("FFmpeg:\n")
			switch runtime.GOOS {
			case "darwin":
				guide.WriteString("  macOS: brew install ffmpeg\n")
			case "windows":
				guide.WriteString("  Windows: choco install ffmpeg\n")
				guide.WriteString("  Or download: https://ffmpeg.org/download.html\n")
			default:
				guide.WriteString("  Linux:\n")
				guide.WriteString("    Ubuntu/Debian: sudo apt install ffmpeg\n")
				guide.WriteString("    CentOS/RHEL: sudo yum install ffmpeg\n")
			}
			guide.WriteString("\n")
		}
	}

	return guide.String()
}
