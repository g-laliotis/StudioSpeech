package agents

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Voice represents a voice model in the catalog
type Voice struct {
	ID                   string `json:"id"`
	Language             string `json:"language"`
	Gender               string `json:"gender"`
	Style                string `json:"style"`
	SampleRate           int    `json:"sample_rate"`
	CommercialUseAllowed bool   `json:"commercial_use_allowed"`
	AttributionRequired  bool   `json:"attribution_required"`
	LicenseName          string `json:"license_name"`
	LicenseURL           string `json:"license_url"`
	SourceURL            string `json:"source_url"`
	SHA256               string `json:"sha256"`
	Path                 string `json:"path"`
}

// VoiceCatalog contains all available voices
type VoiceCatalog struct {
	Voices []Voice `json:"voices"`
}

// VoiceCatalogAgent manages voice model selection and validation
type VoiceCatalogAgent struct {
	catalogPath string
	catalog     *VoiceCatalog
}

// NewVoiceCatalogAgent creates a new voice catalog agent
func NewVoiceCatalogAgent(catalogPath string) *VoiceCatalogAgent {
	return &VoiceCatalogAgent{
		catalogPath: catalogPath,
	}
}

// LoadCatalog reads and validates the voice catalog
func (v *VoiceCatalogAgent) LoadCatalog() error {
	file, err := os.Open(v.catalogPath)
	if err != nil {
		return fmt.Errorf("failed to open catalog file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	catalog := &VoiceCatalog{}

	if err := decoder.Decode(catalog); err != nil {
		return fmt.Errorf("failed to parse catalog JSON: %w", err)
	}

	// Validate each voice entry
	for i, voice := range catalog.Voices {
		if err := v.validateVoice(&voice); err != nil {
			return fmt.Errorf("invalid voice entry %d (%s): %w", i, voice.ID, err)
		}
	}

	v.catalog = catalog
	return nil
}

// validateVoice ensures a voice entry meets commercial safety requirements
func (v *VoiceCatalogAgent) validateVoice(voice *Voice) error {
	// Check required fields
	if voice.ID == "" {
		return fmt.Errorf("voice ID is required")
	}

	if voice.Language == "" {
		return fmt.Errorf("language is required")
	}

	if voice.LicenseName == "" {
		return fmt.Errorf("license_name is required")
	}

	// CRITICAL: Block non-commercial voices
	if !voice.CommercialUseAllowed {
		return fmt.Errorf("voice %s is not allowed for commercial use (license: %s)",
			voice.ID, voice.LicenseName)
	}

	// Auto-detect non-commercial licenses
	licenseLower := strings.ToLower(voice.LicenseName)
	if strings.Contains(licenseLower, "non-commercial") ||
		strings.Contains(licenseLower, "nc") ||
		strings.Contains(licenseLower, "by-nc") {
		return fmt.Errorf("voice %s has non-commercial license: %s",
			voice.ID, voice.LicenseName)
	}

	return nil
}

// SelectVoice chooses appropriate voice based on language, voice ID, and gender
func (v *VoiceCatalogAgent) SelectVoice(language, voiceID, gender string) (*Voice, error) {
	if v.catalog == nil {
		return nil, fmt.Errorf("catalog not loaded")
	}

	// If specific voice ID requested, find it
	if voiceID != "auto" && voiceID != "" {
		for _, voice := range v.catalog.Voices {
			if voice.ID == voiceID {
				return &voice, nil
			}
		}
		return nil, fmt.Errorf("voice ID %s not found in catalog", voiceID)
	}

	// Auto-select based on language and gender
	var candidates []Voice

	// Normalize language code
	lang := v.normalizeLanguage(language)

	// Find voices matching the language
	for _, voice := range v.catalog.Voices {
		if strings.HasPrefix(voice.Language, lang) {
			candidates = append(candidates, voice)
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no voices found for language %s", language)
	}

	// Filter by gender if specified
	if gender != "auto" && gender != "" {
		var genderCandidates []Voice
		for _, voice := range candidates {
			if voice.Gender == gender {
				genderCandidates = append(genderCandidates, voice)
			}
		}
		if len(genderCandidates) > 0 {
			candidates = genderCandidates
		}
	}

	// Prefer higher quality voices (heuristic: higher sample rate)
	bestVoice := &candidates[0]
	for i := 1; i < len(candidates); i++ {
		if candidates[i].SampleRate > bestVoice.SampleRate {
			bestVoice = &candidates[i]
		}
	}

	return bestVoice, nil
}

// normalizeLanguage converts language codes to standard format
func (v *VoiceCatalogAgent) normalizeLanguage(lang string) string {
	switch strings.ToLower(lang) {
	case "en", "english", "en-us", "en_us":
		return "en-US"
	case "en-uk", "en_uk", "en-gb", "en_gb":
		return "en-UK"
	case "el", "greek", "el-gr", "el_gr":
		return "el-GR"
	case "auto":
		return "en-US" // Default to English
	default:
		return lang
	}
}

// ValidateVoiceFile checks if voice model file exists and matches expected hash
func (v *VoiceCatalogAgent) ValidateVoiceFile(voice *Voice) error {
	// Check if file exists
	if _, err := os.Stat(voice.Path); os.IsNotExist(err) {
		return fmt.Errorf("voice model file not found: %s", voice.Path)
	}

	// Skip hash validation if not provided
	if voice.SHA256 == "" || voice.SHA256 == "<fill after download>" {
		return nil
	}

	// Validate file hash
	file, err := os.Open(voice.Path)
	if err != nil {
		return fmt.Errorf("failed to open voice file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to calculate file hash: %w", err)
	}

	actualHash := fmt.Sprintf("%x", hasher.Sum(nil))
	if actualHash != voice.SHA256 {
		return fmt.Errorf("voice file hash mismatch: expected %s, got %s",
			voice.SHA256, actualHash)
	}

	return nil
}

// GetAvailableVoices returns list of all valid voices
func (v *VoiceCatalogAgent) GetAvailableVoices() []Voice {
	if v.catalog == nil {
		return nil
	}
	return v.catalog.Voices
}

// GetAttributionText returns required attribution text for voices that need it
func (v *VoiceCatalogAgent) GetAttributionText() []string {
	var attributions []string

	if v.catalog == nil {
		return attributions
	}

	for _, voice := range v.catalog.Voices {
		if voice.AttributionRequired {
			switch {
			case strings.Contains(strings.ToLower(voice.LicenseName), "libritts"):
				attributions = append(attributions,
					fmt.Sprintf("Voice %s: This project uses the LibriTTS dataset (CC BY 4.0). "+
						"© Original contributors. Licensed under CC BY 4.0 (%s). No endorsement implied.",
						voice.ID, voice.LicenseURL))
			default:
				attributions = append(attributions,
					fmt.Sprintf("Voice %s: %s (%s)", voice.ID, voice.LicenseName, voice.LicenseURL))
			}
		}
	}

	return attributions
}
