package checkstyle

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/statik/pkg/plugin"
)

// CheckstyleFile represents a single file in the Checkstyle output
type CheckstyleFile struct {
	Name    string         `xml:"name,attr"`
	Errors  []CheckstyleError `xml:"error"`
}

// CheckstyleError represents a single error in the Checkstyle output
type CheckstyleError struct {
	Line     int    `xml:"line,attr"`
	Column   int    `xml:"column,attr"`
	Severity string `xml:"severity,attr"`
	Message  string `xml:"message,attr"`
	Source   string `xml:"source,attr"`
}

// CheckstyleOutput represents the root element of the Checkstyle XML output
type CheckstyleOutput struct {
	Files []CheckstyleFile `xml:"file"`
}

// Parser implements the plugin.Parser interface for Checkstyle output
type Parser struct{}

// Name returns the name of the parser
func (p *Parser) Name() string {
	return "checkstyle"
}

// Parse reads Checkstyle XML output and converts it to AnalysisResults
func (p *Parser) Parse(reader io.Reader) ([]plugin.AnalysisResult, error) {
	var output CheckstyleOutput
	if err := xml.NewDecoder(reader).Decode(&output); err != nil {
		return nil, fmt.Errorf("failed to decode Checkstyle output: %w", err)
	}

	results := make([]plugin.AnalysisResult, 0)
	for _, file := range output.Files {
		for _, err := range file.Errors {
			// Convert Checkstyle severity to plugin.Severity
			var severity plugin.Severity
			switch err.Severity {
			case "error":
				severity = plugin.SeverityError
			case "warning":
				severity = plugin.SeverityWarning
			default:
				// Skip other severities like "info"
				continue
			}

			// Extract rule ID from source (e.g., "com.puppycrawl.tools.checkstyle.checks.naming.ConstantNameCheck")
			ruleID := err.Source
			if ruleID == "" {
				ruleID = "checkstyle"
			}

			results = append(results, plugin.AnalysisResult{
				Tool:        p.Name(),
				File:        file.Name,
				Line:        err.Line,
				Column:      err.Column,
				Message:     err.Message,
				Severity:    severity,
				RuleID:      ruleID,
				Description: err.Message,
			})
		}
	}

	return results, nil
}

// SupportedFileExtensions returns the file extensions this parser can handle
func (p *Parser) SupportedFileExtensions() []string {
	return []string{".java", ".xml", ".properties"}
}

// GetRuleSummary returns a custom summary for a specific rule, or nil if no custom summary is needed
func (p *Parser) GetRuleSummary(ruleID string, results []plugin.AnalysisResult) *plugin.RuleSummary {
	// No custom summaries for Checkstyle rules yet
	return nil
} 