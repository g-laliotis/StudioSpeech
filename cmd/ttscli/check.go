package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"studiospeech/internal/agents"
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
	
	// Initialize environment agent
	envAgent := agents.NewEnvironmentAgent()
	envInfo, err := envAgent.Check()
	if err != nil {
		fmt.Printf("❌ Environment check failed: %v\n", err)
		return
	}
	
	// Display system info
	fmt.Printf("Go version: %s\n", envInfo.GoVersion)
	fmt.Printf("OS/Arch: %s/%s\n\n", envInfo.OS, envInfo.Arch)
	
	// Check Piper TTS
	if envInfo.HasPiper {
		fmt.Printf("Checking Piper TTS... ✅ Found: %s\n", envInfo.PiperVersion)
	} else {
		fmt.Println("Checking Piper TTS... ❌ NOT FOUND")
	}
	
	// Check FFmpeg
	if envInfo.HasFFmpeg {
		fmt.Printf("Checking FFmpeg... ✅ Found: FFmpeg %s\n", envInfo.FFmpegVersion)
	} else {
		fmt.Println("Checking FFmpeg... ❌ NOT FOUND")
	}
	
	// Check voice catalog
	checkVoiceCatalog()
	
	// Show installation guide for missing components
	var missing []string
	if !envInfo.HasPiper {
		missing = append(missing, "piper")
	}
	if !envInfo.HasFFmpeg {
		missing = append(missing, "ffmpeg")
	}
	
	if len(missing) > 0 {
		fmt.Println("\n" + envAgent.GetInstallGuide(missing))
	}
	
	fmt.Println("✅ System check complete!")
}

// checkVoiceCatalog validates the voice catalog
func checkVoiceCatalog() {
	fmt.Print("Checking voice catalog... ")
	
	// Initialize voice catalog agent
	catalogPath := filepath.Join("voices", "catalog.json")
	voiceAgent := agents.NewVoiceCatalogAgent(catalogPath)
	
	// Load and validate catalog
	if err := voiceAgent.LoadCatalog(); err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		return
	}
	
	voices := voiceAgent.GetAvailableVoices()
	fmt.Printf("✅ Found %d commercial-safe voices\n", len(voices))
	
	// Display available voices
	for _, voice := range voices {
		fmt.Printf("  - %s (%s, %s %s)\n", voice.ID, voice.Language, voice.Gender, voice.Style)
	}
	
	// Show attribution requirements
	attributions := voiceAgent.GetAttributionText()
	if len(attributions) > 0 {
		fmt.Println("\n📝 Attribution Required:")
		for _, attr := range attributions {
			fmt.Printf("  %s\n", attr)
		}
	}
}