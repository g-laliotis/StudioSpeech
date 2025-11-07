package main

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ttscli",
	Short: "StudioSpeech - Local Text-to-Speech for YouTube",
	Long: `StudioSpeech is a free, local (offline) Text-to-Speech tool for creating 
high-quality speech from text files, designed for commercial YouTube content creation.

Features:
- Local & Offline: No cloud APIs, no per-character fees
- Commercial Safe: Only uses voice models with commercial licensing
- Multi-language: Supports Greek (el-GR) and English (en-US/UK)
- High Quality: Natural-sounding speech suitable for YouTube`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

func main() {
	Execute()
}