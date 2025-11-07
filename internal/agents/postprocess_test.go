package agents

import (
	"strings"
	"testing"
)

func TestPostProcessAgent_ValidateParams(t *testing.T) {
	agent := NewPostProcessAgent("ffmpeg", "/tmp")
	
	tests := []struct {
		params    *PostProcessParams
		shouldErr bool
		desc      string
	}{
		{
			params:    &PostProcessParams{Format: FormatMP3, SampleRate: 48000, Bitrate: 192, LoudnessLUFS: -16.0},
			shouldErr: false,
			desc:      "valid MP3 parameters",
		},
		{
			params:    &PostProcessParams{Format: FormatWAV, SampleRate: 48000, Bitrate: 0, LoudnessLUFS: -16.0},
			shouldErr: false,
			desc:      "valid WAV parameters",
		},
		{
			params:    &PostProcessParams{Format: "invalid", SampleRate: 48000, Bitrate: 192, LoudnessLUFS: -16.0},
			shouldErr: true,
			desc:      "invalid format",
		},
		{
			params:    &PostProcessParams{Format: FormatMP3, SampleRate: 1000, Bitrate: 192, LoudnessLUFS: -16.0},
			shouldErr: true,
			desc:      "sample rate too low",
		},
		{
			params:    &PostProcessParams{Format: FormatMP3, SampleRate: 48000, Bitrate: 32, LoudnessLUFS: -16.0},
			shouldErr: true,
			desc:      "bitrate too low",
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

func TestPostProcessAgent_GetCommandLine(t *testing.T) {
	agent := NewPostProcessAgent("ffmpeg", "/tmp")
	
	params := &PostProcessParams{
		Format:       FormatMP3,
		SampleRate:   48000,
		Bitrate:      192,
		LoudnessLUFS: -16.0,
	}
	
	cmdLine := agent.GetCommandLine("/input.wav", "/output.mp3", params)
	
	expectedParts := []string{
		"-i", "/input.wav",
		"-y",
		"aresample=48000",
		"loudnorm=I=-16.0",
		"-codec:a", "libmp3lame",
		"-b:a", "192k",
		"/output.mp3",
	}
	
	for _, part := range expectedParts {
		if !strings.Contains(cmdLine, part) {
			t.Errorf("Command line missing expected part: %s\nFull command: %s", part, cmdLine)
		}
	}
}

func TestPostProcessAgent_GetDefaultParams(t *testing.T) {
	agent := NewPostProcessAgent("ffmpeg", "/tmp")
	
	params := agent.getDefaultParams()
	
	if params.Format != FormatMP3 {
		t.Errorf("Expected default format MP3, got %s", params.Format)
	}
	
	if params.SampleRate != 48000 {
		t.Errorf("Expected default sample rate 48000, got %d", params.SampleRate)
	}
	
	if params.Bitrate != 192 {
		t.Errorf("Expected default bitrate 192, got %d", params.Bitrate)
	}
}