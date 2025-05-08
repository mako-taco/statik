package plugin

import (
	"io"
)

// Severity represents the severity level of an analysis result
type Severity string

const (
	// SeverityWarning represents a warning severity
	SeverityWarning Severity = "WARNING"
	// SeverityError represents an error severity
	SeverityError Severity = "ERROR"
)

// AnalysisResult represents a single issue found by a static analysis tool
type AnalysisResult struct {
	Tool        string   `json:"tool"`
	File        string   `json:"file"`
	Line        int      `json:"line"`
	Column      int      `json:"column"`
	Message     string   `json:"message"`
	Severity    Severity `json:"severity"`
	RuleID      string   `json:"rule_id"`
	Description string   `json:"description,omitempty"`
}

// Parser defines the interface that all static analysis tool parsers must implement
type Parser interface {
	// Name returns the name of the parser
	Name() string

	// Parse reads from the provided reader and returns a slice of AnalysisResults
	Parse(reader io.Reader) ([]AnalysisResult, error)

	// SupportedFileExtensions returns a list of file extensions this parser can handle
	SupportedFileExtensions() []string
} 