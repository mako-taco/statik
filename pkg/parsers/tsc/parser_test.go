package tsc

import (
	"strings"
	"testing"

	"github.com/statik/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []plugin.AnalysisResult
		wantErr  bool
	}{
		{
			name: "valid error message",
			input: `src/app.ts(10,5): error TS2322: Type 'string' is not assignable to type 'number'.`,
			expected: []plugin.AnalysisResult{
				{
					Tool:        "tsc",
					File:        "src/app.ts",
					Line:        10,
					Column:      5,
					Message:     "Type 'string' is not assignable to type 'number'.",
					Severity:    plugin.SeverityError,
					RuleID:      "TS2322",
					Description: "Type 'string' is not assignable to type 'number'.",
				},
			},
			wantErr: false,
		},
		{
			name: "valid warning message",
			input: `src/utils.ts(15,8): warning TS6133: 'i' is declared but its value is never read.`,
			expected: []plugin.AnalysisResult{
				{
					Tool:        "tsc",
					File:        "src/utils.ts",
					Line:        15,
					Column:      8,
					Message:     "'i' is declared but its value is never read.",
					Severity:    plugin.SeverityWarning,
					RuleID:      "TS6133",
					Description: "'i' is declared but its value is never read.",
				},
			},
			wantErr: false,
		},
		{
			name:     "invalid message format",
			input:    "This is not a valid TypeScript error message",
			expected: []plugin.AnalysisResult{},
			wantErr:  false,
		},
		{
			name: "multiple messages",
			input: `src/app.ts(10,5): error TS2322: Type 'string' is not assignable to type 'number'.
src/utils.ts(15,8): warning TS6133: 'i' is declared but its value is never read.`,
			expected: []plugin.AnalysisResult{
				{
					Tool:        "tsc",
					File:        "src/app.ts",
					Line:        10,
					Column:      5,
					Message:     "Type 'string' is not assignable to type 'number'.",
					Severity:    plugin.SeverityError,
					RuleID:      "TS2322",
					Description: "Type 'string' is not assignable to type 'number'.",
				},
				{
					Tool:        "tsc",
					File:        "src/utils.ts",
					Line:        15,
					Column:      8,
					Message:     "'i' is declared but its value is never read.",
					Severity:    plugin.SeverityWarning,
					RuleID:      "TS6133",
					Description: "'i' is declared but its value is never read.",
				},
			},
			wantErr: false,
		},
	}

	parser := &Parser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			results, err := parser.Parse(reader)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, results)
		})
	}
}

func TestParser_Name(t *testing.T) {
	parser := &Parser{}
	assert.Equal(t, "tsc", parser.Name())
}

func TestParser_SupportedFileExtensions(t *testing.T) {
	parser := &Parser{}
	extensions := parser.SupportedFileExtensions()
	assert.Contains(t, extensions, ".ts")
	assert.Contains(t, extensions, ".tsx")
	assert.Len(t, extensions, 2)
} 