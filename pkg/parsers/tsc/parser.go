package tsc

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/statik/pkg/plugin"
)

// Parser implements the plugin.Parser interface for TypeScript compiler output
type Parser struct{}

// Name returns the name of the parser
func (p *Parser) Name() string {
	return "tsc"
}

// Parse reads TypeScript compiler output and converts it to AnalysisResults
func (p *Parser) Parse(reader io.Reader) ([]plugin.AnalysisResult, error) {
	results := make([]plugin.AnalysisResult, 0)
	scanner := bufio.NewScanner(reader)

	// Regular expression to match TypeScript error lines
	// Format: file.ts(line,col): error TS1234: Error message
	re := regexp.MustCompile(`^(.+)\((\d+),(\d+)\): (error|warning) (TS\d+): (.+)$`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) != 7 {
			continue
		}

		file := matches[1]
		lineNum, _ := strconv.Atoi(matches[2])
		colNum, _ := strconv.Atoi(matches[3])
		severityStr := matches[4]
		ruleID := matches[5]
		message := matches[6]

		// Convert severity string to Severity type
		var severity plugin.Severity
		if severityStr == "error" {
			severity = plugin.SeverityError
		} else {
			severity = plugin.SeverityWarning
		}

		results = append(results, plugin.AnalysisResult{
			Tool:        p.Name(),
			File:        file,
			Line:        lineNum,
			Column:      colNum,
			Message:     message,
			Severity:    severity,
			RuleID:      ruleID,
			Description: message,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading TypeScript output: %w", err)
	}

	return results, nil
}

// SupportedFileExtensions returns the file extensions this parser can handle
func (p *Parser) SupportedFileExtensions() []string {
	return []string{".ts", ".tsx"}
} 