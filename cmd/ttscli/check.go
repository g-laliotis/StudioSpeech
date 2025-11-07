package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check system requirements and voice availability",
	Long: `Check if all required dependencies are installed and available:
- Piper TTS (for speech synthesis)
- FFmpeg (for audio processing)
- Available voice models

This command helps diagnose setup issues and provides installation guidance.`,
	Run: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

// runCheck performs system requirements validation
func runCheck(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 StudioSpeech System Check")
	fmt.Println("============================")
	
	// Check Go version
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
	
	// Check Piper TTS
	checkPiper()
	
	// Check FFmpeg
	checkFFmpeg()
	
	// Check voice catalog
	checkVoices()
	
	fmt.Println("\n✅ System check complete!")
}

// checkPiper verifies Piper TTS installation
func checkPiper() {
	fmt.Print("Checking Piper TTS... ")
	
	cmd := exec.Command("piper", "--version")
	output, err := cmd.Output()
	
	if err != nil {
		fmt.Println("❌ NOT FOUND")
		printPiperInstallGuide()
		return
	}
	
	version := strings.TrimSpace(string(output))
	fmt.Printf("✅ Found: %s\n", version)
}

// checkFFmpeg verifies FFmpeg installation
func checkFFmpeg() {
	fmt.Print("Checking FFmpeg... ")
	
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.Output()
	
	if err != nil {
		fmt.Println("❌ NOT FOUND")
		printFFmpegInstallGuide()
		return
	}
	
	// Extract version from first line
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		version := strings.Fields(lines[0])
		if len(version) >= 3 {
			fmt.Printf("✅ Found: %s %s\n", version[0], version[2])
		} else {
			fmt.Println("✅ Found: FFmpeg (version unknown)")
		}
	}
}

// checkVoices validates voice catalog
func checkVoices() {
	fmt.Print("Checking voice catalog... ")
	// TODO: Implement voice catalog validation
	fmt.Println("⚠️  Voice catalog validation not yet implemented")
}

// printPiperInstallGuide provides OS-specific installation instructions for Piper
func printPiperInstallGuide() {
	fmt.Println("\n📋 Piper TTS Installation Guide:")
	
	switch runtime.GOOS {
	case "darwin": // macOS
		fmt.Println("  macOS:")
		fmt.Println("    brew install piper-tts")
		fmt.Println("  Or download from: https://github.com/rhasspy/piper/releases")
		
	case "windows":
		fmt.Println("  Windows:")
		fmt.Println("    choco install piper-tts")
		fmt.Println("  Or download from: https://github.com/rhasspy/piper/releases")
		
	default:
		fmt.Println("  Linux:")
		fmt.Println("    Download from: https://github.com/rhasspy/piper/releases")
		fmt.Println("    Extract and add to PATH")
	}
}

// printFFmpegInstallGuide provides OS-specific installation instructions for FFmpeg
func printFFmpegInstallGuide() {
	fmt.Println("\n📋 FFmpeg Installation Guide:")
	
	switch runtime.GOOS {
	case "darwin": // macOS
		fmt.Println("  macOS:")
		fmt.Println("    brew install ffmpeg")
		
	case "windows":
		fmt.Println("  Windows:")
		fmt.Println("    choco install ffmpeg")
		fmt.Println("  Or download from: https://ffmpeg.org/download.html")
		
	default:
		fmt.Println("  Linux:")
		fmt.Println("    sudo apt install ffmpeg  # Ubuntu/Debian")
		fmt.Println("    sudo yum install ffmpeg  # CentOS/RHEL")
	}
}