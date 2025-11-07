package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// synthCmd represents the synth command
var synthCmd = &cobra.Command{
	Use:   "synth",
	Short: "Synthesize speech from text files",
	Long: `Convert text files (.txt, .docx, or .pdf) to high-quality speech audio.

Examples:
  ttscli synth --in script.txt --lang en-US --gender female --out voice.mp3
  ttscli synth --in document.docx --lang el-GR --gender male --format wav
  ttscli synth --in document.pdf --lang en-US --gender female --out speech.mp3
  ttscli synth --in story.txt --speed 1.05 --gender auto --out narration.mp3`,
	Run: runSynth,
}

var (
	// Input/Output flags
	inputFile  string
	outputFile string
	
	// Language and voice flags
	language string
	voiceID  string
	gender   string
	
	// Audio format flags
	format     string
	sampleRate int
	bitrate    int
	
	// Synthesis parameters
	speed  float64
	noise  float64
	noisew float64
	
	// Processing flags
	noCache bool
)

func init() {
	rootCmd.AddCommand(synthCmd)
	
	// Input/Output flags
	synthCmd.Flags().StringVarP(&inputFile, "in", "i", "", "input text file (.txt or .docx)")
	synthCmd.Flags().StringVarP(&outputFile, "out", "o", "", "output audio file")
	synthCmd.MarkFlagRequired("in")
	synthCmd.MarkFlagRequired("out")
	
	// Language and voice flags
	synthCmd.Flags().StringVarP(&language, "lang", "l", "auto", "language code (en-US, en-UK, el-GR, or auto)")
	synthCmd.Flags().StringVar(&voiceID, "voice", "auto", "voice ID from catalog (or auto for default)")
	synthCmd.Flags().StringVarP(&gender, "gender", "g", "auto", "voice gender (male, female, or auto)")

	
	// Audio format flags
	synthCmd.Flags().StringVarP(&format, "format", "f", "mp3", "output format (wav, mp3)")
	synthCmd.Flags().IntVar(&sampleRate, "sample-rate", 48000, "output sample rate in Hz")
	synthCmd.Flags().IntVar(&bitrate, "bitrate", 192, "MP3 bitrate in kbps")
	
	// Synthesis parameters
	synthCmd.Flags().Float64Var(&speed, "speed", 1.03, "speech speed multiplier (0.5-2.0)")
	synthCmd.Flags().Float64Var(&noise, "noise", 0.667, "noise level for naturalness (0.0-1.0)")
	synthCmd.Flags().Float64Var(&noisew, "noisew", 0.8, "noise width for variation (0.0-1.0)")
	
	// Processing flags
	synthCmd.Flags().BoolVar(&noCache, "no-cache", false, "disable caching of synthesized audio")
}

// runSynth executes the speech synthesis pipeline
func runSynth(cmd *cobra.Command, args []string) {
	fmt.Printf("🎤 StudioSpeech Synthesis\n")
	fmt.Printf("========================\n\n")
	
	// Validate input file
	if err := validateInputFile(inputFile); err != nil {
		fmt.Printf("❌ Input file error: %v\n", err)
		return
	}
	
	// Validate output file
	if err := validateOutputFile(outputFile); err != nil {
		fmt.Printf("❌ Output file error: %v\n", err)
		return
	}
	
	// Validate parameters
	if err := validateSynthParams(); err != nil {
		fmt.Printf("❌ Parameter error: %v\n", err)
		return
	}
	
	fmt.Printf("📄 Input: %s\n", inputFile)
	fmt.Printf("🔊 Output: %s\n", outputFile)
	fmt.Printf("🌍 Language: %s\n", language)
	fmt.Printf("🎭 Voice: %s\n", voiceID)
	fmt.Printf("⚡ Speed: %.2fx\n", speed)
	fmt.Printf("📊 Format: %s\n", strings.ToUpper(format))
	
	if format == "mp3" {
		fmt.Printf("🎵 Bitrate: %d kbps\n", bitrate)
	}
	fmt.Printf("📈 Sample Rate: %d Hz\n\n", sampleRate)
	
	// Execute synthesis pipeline
	if err := executeSynthesisPipeline(); err != nil {
		fmt.Printf("❌ Synthesis failed: %v\n", err)
		return
	}
	
	fmt.Println("✅ Synthesis completed successfully!")
}

// validateInputFile checks if input file exists and has supported extension
func validateInputFile(path string) error {
	if path == "" {
		return fmt.Errorf("input file is required")
	}
	
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".txt" && ext != ".docx" && ext != ".pdf" {
		return fmt.Errorf("unsupported file type: %s (supported: .txt, .docx, .pdf)", ext)
	}
	
	// TODO: Check if file exists
	return nil
}

// validateOutputFile checks output file path and extension
func validateOutputFile(path string) error {
	if path == "" {
		return fmt.Errorf("output file is required")
	}
	
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".wav" && ext != ".mp3" {
		return fmt.Errorf("unsupported output format: %s (supported: .wav, .mp3)", ext)
	}
	
	// Auto-detect format from extension if not explicitly set
	if format == "mp3" && ext == ".wav" {
		format = "wav"
	} else if format == "wav" && ext == ".mp3" {
		format = "mp3"
	}
	
	return nil
}

// validateSynthParams validates synthesis parameters
func validateSynthParams() error {
	if speed < 0.5 || speed > 2.0 {
		return fmt.Errorf("speed must be between 0.5 and 2.0, got %.2f", speed)
	}
	
	if noise < 0.0 || noise > 1.0 {
		return fmt.Errorf("noise must be between 0.0 and 1.0, got %.3f", noise)
	}
	
	if noisew < 0.0 || noisew > 1.0 {
		return fmt.Errorf("noisew must be between 0.0 and 1.0, got %.3f", noisew)
	}
	
	if bitrate < 64 || bitrate > 320 {
		return fmt.Errorf("bitrate must be between 64 and 320 kbps, got %d", bitrate)
	}
	
	return nil
}