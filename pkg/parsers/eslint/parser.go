package eslint

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/statik/pkg/plugin"
)

// ESLintMessage represents a single ESLint message
type ESLintMessage struct {
	RuleID    string `json:"ruleId"`
	Severity  int    `json:"severity"`
	Message   string `json:"message"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	NodeType  string `json:"nodeType,omitempty"`
	MessageID string `json:"messageId,omitempty"`
	EndLine   int    `json:"endLine,omitempty"`
	EndColumn int    `json:"endColumn,omitempty"`
	Fix       struct {
		Range []int  `json:"range"`
		Text  string `json:"text"`
	} `json:"fix,omitempty"`
}

// ESLintFile represents ESLint output for a single file
type ESLintFile struct {
	FilePath string         `json:"filePath"`
	Messages []ESLintMessage `json:"messages"`
	ErrorCount int          `json:"errorCount"`
	WarningCount int        `json:"warningCount"`
	FixableErrorCount int   `json:"fixableErrorCount"`
	FixableWarningCount int `json:"fixableWarningCount"`
	Source string          `json:"source,omitempty"`
}

// Parser implements the plugin.Parser interface for ESLint output
type Parser struct{}

// Name returns the name of the parser
func (p *Parser) Name() string {
	return "eslint"
}

// Parse reads ESLint JSON output and converts it to AnalysisResults
func (p *Parser) Parse(reader io.Reader) ([]plugin.AnalysisResult, error) {
	var files []ESLintFile
	if err := json.NewDecoder(reader).Decode(&files); err != nil {
		return nil, fmt.Errorf("failed to decode ESLint output: %w", err)
	}

	results := make([]plugin.AnalysisResult, 0)
	for _, file := range files {
		for _, msg := range file.Messages {
			// Skip messages without a rule ID (usually parsing errors)
			if msg.RuleID == "" {
				continue
			}

			// Convert ESLint severity (2=error, 1=warning, 0=off) to Severity
			severity := plugin.SeverityWarning
			if msg.Severity == 2 {
				severity = plugin.SeverityError
			}

			results = append(results, plugin.AnalysisResult{
				Tool:        p.Name(),
				File:        file.FilePath,
				Line:        msg.Line,
				Column:      msg.Column,
				Message:     msg.Message,
				Severity:    severity,
				RuleID:      msg.RuleID,
				Description: msg.Message,
			})
		}
	}

	return results, nil
}

// SupportedFileExtensions returns the file extensions this parser can handle
func (p *Parser) SupportedFileExtensions() []string {
	return []string{".js", ".jsx", ".ts", ".tsx", ".vue"}
}

// GetRuleSummary returns a custom summary for a specific rule, or nil if no custom summary is needed
func (p *Parser) GetRuleSummary(ruleID string, results []plugin.AnalysisResult) *plugin.RuleSummary {
	if ruleID != "max-lines" {
		return nil
	}

	if len(results) == 0 {
		return nil
	}

	// For max-lines, we want to show the difference between current and max lines
	// The message format is typically: "File has too many lines (X). Maximum allowed is Y."
	// We'll extract these numbers and use the difference as the count
	message := results[0].Message
	var currentLines, maxLines int
	fmt.Sscanf(message, "File has too many lines (%d). Maximum allowed is %d.", &currentLines, &maxLines)

	// Create a custom summary
	summary := &plugin.RuleSummary{
		RuleID:      ruleID,
		Description: results[0].Description,
		Severity:    results[0].Severity,
		Count:       currentLines - maxLines, // Use the difference as the count
		Violations:  make([]plugin.Violation, 0, len(results)),
	}

	// Add violations
	for _, result := range results {
		summary.Violations = append(summary.Violations, plugin.Violation{
			Line:    result.Line,
			Column:  result.Column,
			Message: result.Message,
		})
	}

	return summary
} 