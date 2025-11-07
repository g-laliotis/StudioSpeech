package agents

import (
	"strings"
	"testing"
)

func TestNormalizeAgent_CleanText(t *testing.T) {
	agent := NewNormalizeAgent()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello world and test",
			expected: "Hello world and test",
		},
		{
			input:    "Text with — em dash and – en dash",
			expected: "Text with - em dash and - en dash",
		},
		{
			input:    "Multiple    spaces   here",
			expected: "Multiple spaces here",
		},
	}

	for _, test := range tests {
		result := agent.cleanText(test.input)
		if result != test.expected {
			t.Errorf("cleanText(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestNormalizeAgent_ExpandAbbreviations(t *testing.T) {
	agent := NewNormalizeAgent()

	tests := []struct {
		input    string
		language string
		expected string
	}{
		{
			input:    "Dr. Smith went to the store etc.",
			language: "en-US",
			expected: "Doctor Smith went to the store etcetera",
		},
		{
			input:    "Το κείμενο κ.λπ. είναι εδώ.",
			language: "el-GR",
			expected: "Το κείμενο και λοιπά είναι εδώ.",
		},
	}

	for _, test := range tests {
		result := agent.expandAbbreviations(test.input, test.language)
		if result != test.expected {
			t.Errorf("expandAbbreviations(%q, %q) = %q, want %q",
				test.input, test.language, result, test.expected)
		}
	}
}

func TestNormalizeAgent_ExpandNumbers(t *testing.T) {
	agent := NewNormalizeAgent()

	tests := []struct {
		input    string
		language string
		expected string
	}{
		{
			input:    "I have 5 apples and 10 oranges.",
			language: "en-US",
			expected: "I have five apples and ten oranges.",
		},
		{
			input:    "Έχω 3 μήλα και 7 πορτοκάλια.",
			language: "el-GR",
			expected: "Έχω τρία μήλα και επτά πορτοκάλια.",
		},
	}

	for _, test := range tests {
		result := agent.expandNumbers(test.input, test.language)
		if result != test.expected {
			t.Errorf("expandNumbers(%q, %q) = %q, want %q",
				test.input, test.language, result, test.expected)
		}
	}
}

func TestNormalizeAgent_Normalize(t *testing.T) {
	agent := NewNormalizeAgent()

	content := &TextContent{
		Paragraphs: []string{
			"Hello Dr. Smith. I have 5 apples.",
			"This is the 2nd paragraph etc.",
		},
		Language:  "en-US",
		WordCount: 12,
	}

	result, err := agent.Normalize(content)
	if err != nil {
		t.Fatalf("Normalize failed: %v", err)
	}

	if len(result.Sentences) == 0 {
		t.Error("No sentences produced")
	}

	if result.Language != "en-US" {
		t.Errorf("Expected language en-US, got %s", result.Language)
	}

	// Check that abbreviations and numbers were expanded
	allText := strings.Join(result.Sentences, " ")
	if !strings.Contains(allText, "Doctor") {
		t.Error("Dr. was not expanded to Doctor")
	}

	if !strings.Contains(allText, "five") {
		t.Error("5 was not expanded to five")
	}
}
