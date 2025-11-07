package agents

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSynthAgent_ValidateParams(t *testing.T) {
	agent := NewSynthAgent("piper", "/tmp")

	tests := []struct {
		params    *SynthParams
		shouldErr bool
		desc      string
	}{
		{
			params:    &SynthParams{Speed: 1.0, Noise: 0.5, NoiseW: 0.5, Speaker: 0},
			shouldErr: false,
			desc:      "valid parameters",
		},
		{
			params:    &SynthParams{Speed: 0.3, Noise: 0.5, NoiseW: 0.5, Speaker: 0},
			shouldErr: true,
			desc:      "speed too low",
		},
		{
			params:    &SynthParams{Speed: 2.5, Noise: 0.5, NoiseW: 0.5, Speaker: 0},
			shouldErr: true,
			desc:      "speed too high",
		},
		{
			params:    &SynthParams{Speed: 1.0, Noise: -0.1, NoiseW: 0.5, Speaker: 0},
			shouldErr: true,
			desc:      "noise too low",
		},
		{
			params:    &SynthParams{Speed: 1.0, Noise: 1.1, NoiseW: 0.5, Speaker: 0},
			shouldErr: true,
			desc:      "noise too high",
		},
	}

	for _, test := range tests {
		err := agent.validateParams(test.params)
		if test.shouldErr && err == nil {
			t.Errorf("%s: expected error but got none", test.desc)
		}
		if !test.shouldErr && err != nil {
			t.Errorf("%s: unexpected error: %v", test.desc, err)
		}
	}
}

func TestSynthAgent_GetDefaultParams(t *testing.T) {
	agent := NewSynthAgent("piper", "/tmp")

	params := agent.getDefaultParams()

	if params.Speed <= 0 {
		t.Error("Default speed should be positive")
	}

	if params.Noise < 0 || params.Noise > 1 {
		t.Error("Default noise should be between 0 and 1")
	}

	if params.NoiseW < 0 || params.NoiseW > 1 {
		t.Error("Default noiseW should be between 0 and 1")
	}
}

func TestSynthAgent_GetCommandLine(t *testing.T) {
	agent := NewSynthAgent("piper", "/tmp")

	voice := &Voice{
		ID:         "test_voice",
		Path:       "/path/to/voice.onnx",
		SampleRate: 22050,
	}

	params := &SynthParams{
		Speed:   1.0,
		Noise:   0.5,
		NoiseW:  0.8,
		Speaker: 0,
	}

	cmdLine := agent.GetCommandLine(voice, params, "/tmp/output.wav")

	// Check that command contains expected parameters
	expectedParts := []string{
		"--model", "/path/to/voice.onnx",
		"--output_file", "/tmp/output.wav",
		"--length_scale", "1.000",
		"--noise_scale", "0.500",
		"--noise_w", "0.800",
	}

	for _, part := range expectedParts {
		if !strings.Contains(cmdLine, part) {
			t.Errorf("Command line missing expected part: %s\nFull command: %s", part, cmdLine)
		}
	}
}

func TestSynthAgent_DryRun(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "synth_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	agent := NewSynthAgent("piper", tempDir)
	agent.SetDryRun(true)

	normalized := &NormalizedText{
		Sentences: []string{"Hello world.", "This is a test."},
		Language:  "en-US",
	}

	voice := &Voice{
		ID:         "test_voice",
		Path:       filepath.Join(tempDir, "voice.onnx"), // Non-existent file for dry run
		SampleRate: 22050,
	}

	// Create dummy voice file for dry run test
	if err := os.WriteFile(voice.Path, []byte("dummy"), 0644); err != nil {
		t.Fatalf("Failed to create dummy voice file: %v", err)
	}

	params := &SynthParams{
		Speed:   1.0,
		Noise:   0.5,
		NoiseW:  0.8,
		Speaker: 0,
	}

	result, err := agent.Synthesize(normalized, voice, params)
	if err != nil {
		t.Fatalf("Dry run synthesis failed: %v", err)
	}

	if result.OutputPath == "" {
		t.Error("Expected output path in dry run result")
	}

	if result.SampleRate != voice.SampleRate {
		t.Errorf("Expected sample rate %d, got %d", voice.SampleRate, result.SampleRate)
	}
}
