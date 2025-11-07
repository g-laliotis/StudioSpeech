package agents

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// NormalizedText represents processed text ready for synthesis
type NormalizedText struct {
	Sentences []string
	Language  string
	Metadata  map[string]interface{}
}

// NormalizeAgent handles text cleanup and prosody preparation
type NormalizeAgent struct {
	englishAbbrevs map[string]string
	greekAbbrevs   map[string]string
	englishNumbers map[string]string
	greekNumbers   map[string]string
}

// NewNormalizeAgent creates a new text normalization agent
func NewNormalizeAgent() *NormalizeAgent {
	return &NormalizeAgent{
		englishAbbrevs: map[string]string{
			"Dr.":   "Doctor",
			"Mr.":   "Mister",
			"Mrs.":  "Missus",
			"Ms.":   "Miss",
			"Prof.": "Professor",
			"etc.":  "etcetera",
			"vs.":   "versus",
			"e.g.":  "for example",
			"i.e.":  "that is",
		},
		greekAbbrevs: map[string]string{
			"κ.λπ.": "και λοιπά",
			"κ.ά.":  "και άλλα",
			"π.χ.":  "παραδείγματος χάρη",
			"δηλ.":  "δηλαδή",
			"κτλ.":  "και τα λοιπά",
		},
		englishNumbers: map[string]string{
			"0": "zero", "1": "one", "2": "two", "3": "three", "4": "four",
			"5": "five", "6": "six", "7": "seven", "8": "eight", "9": "nine",
			"10": "ten", "11": "eleven", "12": "twelve", "13": "thirteen",
			"14": "fourteen", "15": "fifteen", "16": "sixteen", "17": "seventeen",
			"18": "eighteen", "19": "nineteen", "20": "twenty",
		},
		greekNumbers: map[string]string{
			"0": "μηδέν", "1": "ένα", "2": "δύο", "3": "τρία", "4": "τέσσερα",
			"5": "πέντε", "6": "έξι", "7": "επτά", "8": "οκτώ", "9": "εννέα",
			"10": "δέκα", "11": "έντεκα", "12": "δώδεκα", "13": "δεκατρία",
			"14": "δεκατέσσερα", "15": "δεκαπέντε", "16": "δεκαέξι",
			"17": "δεκαεπτά", "18": "δεκαοκτώ", "19": "δεκαεννέα", "20": "είκοσι",
		},
	}
}

// Normalize processes text content and prepares it for synthesis
func (n *NormalizeAgent) Normalize(content *TextContent) (*NormalizedText, error) {
	if content == nil {
		return nil, fmt.Errorf("content is nil")
	}

	var allSentences []string
	
	for _, paragraph := range content.Paragraphs {
		// Clean and normalize the paragraph
		cleaned := n.cleanText(paragraph)
		
		// Expand abbreviations based on language
		expanded := n.expandAbbreviations(cleaned, content.Language)
		
		// Expand numbers to words
		withNumbers := n.expandNumbers(expanded, content.Language)
		
		// Split into sentences
		sentences := n.splitIntoSentences(withNumbers)
		
		allSentences = append(allSentences, sentences...)
	}

	return &NormalizedText{
		Sentences: allSentences,
		Language:  content.Language,
		Metadata: map[string]interface{}{
			"original_paragraphs": len(content.Paragraphs),
			"total_sentences":     len(allSentences),
			"word_count":         content.WordCount,
		},
	}, nil
}

// cleanText performs basic text cleanup
func (n *NormalizeAgent) cleanText(text string) string {
	// Normalize dashes
	text = strings.ReplaceAll(text, "—", " - ")
	text = strings.ReplaceAll(text, "–", " - ")
	
	// Normalize multiple spaces
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")
	
	return strings.TrimSpace(text)
}

// expandAbbreviations replaces common abbreviations with full words
func (n *NormalizeAgent) expandAbbreviations(text, language string) string {
	var abbrevs map[string]string
	
	switch language {
	case "el-GR":
		abbrevs = n.greekAbbrevs
	default:
		abbrevs = n.englishAbbrevs
	}
	
	for abbrev, expansion := range abbrevs {
		// Simple string replacement for abbreviations
		text = strings.ReplaceAll(text, abbrev, expansion)
	}
	
	return text
}

// expandNumbers converts digits to words for better pronunciation
func (n *NormalizeAgent) expandNumbers(text, language string) string {
	var numbers map[string]string
	
	switch language {
	case "el-GR":
		numbers = n.greekNumbers
	default:
		numbers = n.englishNumbers
	}
	
	// Find standalone numbers (not part of larger numbers or dates)
	numberRegex := regexp.MustCompile(`\b(\d{1,2})\b`)
	
	text = numberRegex.ReplaceAllStringFunc(text, func(match string) string {
		num := strings.TrimSpace(match)
		if expansion, exists := numbers[num]; exists {
			return expansion
		}
		return match // Keep original if no expansion found
	})
	
	return text
}

// splitIntoSentences breaks text into individual sentences
func (n *NormalizeAgent) splitIntoSentences(text string) []string {
	// Simple sentence splitting on common punctuation
	sentenceRegex := regexp.MustCompile(`[.!?]+\s+`)
	
	// Split and clean up
	parts := sentenceRegex.Split(text, -1)
	var sentences []string
	
	for _, part := range parts {
		sentence := strings.TrimSpace(part)
		if sentence != "" {
			// Ensure sentence ends with punctuation
			if !strings.HasSuffix(sentence, ".") && 
			   !strings.HasSuffix(sentence, "!") && 
			   !strings.HasSuffix(sentence, "?") {
				sentence += "."
			}
			sentences = append(sentences, sentence)
		}
	}
	
	return sentences
}

// ProcessPauseMarkup handles optional pause markup like [PAUSE=300ms]
func (n *NormalizeAgent) ProcessPauseMarkup(text string) string {
	// Convert pause markup to sentence breaks for Piper
	pauseRegex := regexp.MustCompile(`\[PAUSE=(\d+)ms\]`)
	
	return pauseRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract duration
		matches := pauseRegex.FindStringSubmatch(match)
		if len(matches) > 1 {
			if duration, err := strconv.Atoi(matches[1]); err == nil {
				// Convert to appropriate punctuation based on duration
				if duration >= 500 {
					return ". " // Long pause - sentence break
				} else if duration >= 200 {
					return ", " // Medium pause - comma
				}
			}
		}
		return " " // Short pause - space
	})
}

// ValidateNormalizedText checks if normalized text is ready for synthesis
func (n *NormalizeAgent) ValidateNormalizedText(normalized *NormalizedText) error {
	if normalized == nil {
		return fmt.Errorf("normalized text is nil")
	}
	
	if len(normalized.Sentences) == 0 {
		return fmt.Errorf("no sentences found after normalization")
	}
	
	// Check for reasonable sentence lengths
	for i, sentence := range normalized.Sentences {
		if len(sentence) > 500 {
			return fmt.Errorf("sentence %d too long: %d characters (max 500)", i, len(sentence))
		}
		
		if strings.TrimSpace(sentence) == "" {
			return fmt.Errorf("sentence %d is empty", i)
		}
	}
	
	return nil
}