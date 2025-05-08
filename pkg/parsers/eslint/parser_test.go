package eslint

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
			name: "valid ESLint output with error and warning",
			input: `[
				{
					"filePath": "src/app.js",
					"messages": [
						{
							"ruleId": "no-unused-vars",
							"severity": 2,
							"message": "'x' is defined but never used",
							"line": 10,
							"column": 5,
							"nodeType": "Identifier"
						},
						{
							"ruleId": "semi",
							"severity": 1,
							"message": "Missing semicolon",
							"line": 15,
							"column": 20,
							"nodeType": "ExpressionStatement"
						}
					],
					"errorCount": 1,
					"warningCount": 1,
					"fixableErrorCount": 0,
					"fixableWarningCount": 1
				}
			]`,
			expected: []plugin.AnalysisResult{
				{
					Tool:        "eslint",
					File:        "src/app.js",
					Line:        10,
					Column:      5,
					Message:     "'x' is defined but never used",
					Severity:    plugin.SeverityError,
					RuleID:      "no-unused-vars",
					Description: "'x' is defined but never used",
				},
				{
					Tool:        "eslint",
					File:        "src/app.js",
					Line:        15,
					Column:      20,
					Message:     "Missing semicolon",
					Severity:    plugin.SeverityWarning,
					RuleID:      "semi",
					Description: "Missing semicolon",
				},
			},
			wantErr: false,
		},
		{
			name: "message without rule ID is skipped",
			input: `[
				{
					"filePath": "src/app.js",
					"messages": [
						{
							"severity": 2,
							"message": "Parsing error: Unexpected token",
							"line": 1,
							"column": 1
						}
					],
					"errorCount": 1,
					"warningCount": 0,
					"fixableErrorCount": 0,
					"fixableWarningCount": 0
				}
			]`,
			expected: []plugin.AnalysisResult{},
			wantErr: false,
		},
		{
			name:     "invalid JSON",
			input:    "invalid json",
			expected: nil,
			wantErr:  true,
		},
		{
			name: "multiple files",
			input: `[
				{
					"filePath": "src/app.js",
					"messages": [
						{
							"ruleId": "no-unused-vars",
							"severity": 2,
							"message": "'x' is defined but never used",
							"line": 10,
							"column": 5
						}
					],
					"errorCount": 1,
					"warningCount": 0,
					"fixableErrorCount": 0,
					"fixableWarningCount": 0
				},
				{
					"filePath": "src/utils.js",
					"messages": [
						{
							"ruleId": "semi",
							"severity": 1,
							"message": "Missing semicolon",
							"line": 15,
							"column": 20
						}
					],
					"errorCount": 0,
					"warningCount": 1,
					"fixableErrorCount": 0,
					"fixableWarningCount": 1
				}
			]`,
			expected: []plugin.AnalysisResult{
				{
					Tool:        "eslint",
					File:        "src/app.js",
					Line:        10,
					Column:      5,
					Message:     "'x' is defined but never used",
					Severity:    plugin.SeverityError,
					RuleID:      "no-unused-vars",
					Description: "'x' is defined but never used",
				},
				{
					Tool:        "eslint",
					File:        "src/utils.js",
					Line:        15,
					Column:      20,
					Message:     "Missing semicolon",
					Severity:    plugin.SeverityWarning,
					RuleID:      "semi",
					Description: "Missing semicolon",
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
	assert.Equal(t, "eslint", parser.Name())
}

func TestParser_SupportedFileExtensions(t *testing.T) {
	parser := &Parser{}
	extensions := parser.SupportedFileExtensions()
	expected := []string{".js", ".jsx", ".ts", ".tsx", ".vue"}
	assert.ElementsMatch(t, expected, extensions)
} 