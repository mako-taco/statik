package plugin

// ToolSummary represents the top-level summary of analysis results
type ToolSummary struct {
	Tool          string         `json:"tool"`
	FileSummaries []FileSummary  `json:"file_summaries"`
}

// FileSummary represents a summary of issues found in a specific file
type FileSummary struct {
	File          string         `json:"file"`
	RuleSummaries []RuleSummary  `json:"rule_summaries"`
}

// RuleSummary represents a summary of issues for a specific rule
type RuleSummary struct {
	RuleID      string     `json:"rule_id"`
	Description string     `json:"description,omitempty"`
	Severity    Severity   `json:"severity"`
	Count       int        `json:"count"`
	Violations  []Violation `json:"violations"`
}

// Violation represents a single instance of a rule violation
type Violation struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
}

// ComparisonResult represents the result of comparing two ToolSummaries
type ComparisonResult struct {
	ImprovedFiles []FileComparison `json:"improved_files"`
	WorsenedFiles []FileComparison `json:"worsened_files"`
	NewFiles      []string         `json:"new_files"`
	RemovedFiles  []string         `json:"removed_files"`
}

// FileComparison represents the comparison of a single file between two summaries
type FileComparison struct {
	File           string          `json:"file"`
	ImprovedRules  []RuleComparison `json:"improved_rules"`
	WorsenedRules  []RuleComparison `json:"worsened_rules"`
	NewRules       []string         `json:"new_rules"`
	RemovedRules   []string         `json:"removed_rules"`
	TotalBefore    int             `json:"total_before"`
	TotalAfter     int             `json:"total_after"`
	NetChange      int             `json:"net_change"`
}

// RuleComparison represents the comparison of a single rule between two summaries
type RuleComparison struct {
	RuleID      string   `json:"rule_id"`
	CountBefore int      `json:"count_before"`
	CountAfter  int      `json:"count_after"`
	Change      int      `json:"change"`
	Severity    Severity `json:"severity"`
}

// Compare compares two ToolSummaries and returns a ComparisonResult
func (s *ToolSummary) Compare(other *ToolSummary) *ComparisonResult {
	if s.Tool != other.Tool {
		return nil // Can't compare different tools
	}

	result := &ComparisonResult{
		ImprovedFiles: make([]FileComparison, 0),
		WorsenedFiles: make([]FileComparison, 0),
		NewFiles:      make([]string, 0),
		RemovedFiles:  make([]string, 0),
	}

	// Create maps for easier lookup
	beforeFiles := make(map[string]FileSummary)
	afterFiles := make(map[string]FileSummary)

	for _, fs := range s.FileSummaries {
		beforeFiles[fs.File] = fs
	}
	for _, fs := range other.FileSummaries {
		afterFiles[fs.File] = fs
	}

	// Compare files that exist in both summaries
	for file, beforeFS := range beforeFiles {
		if afterFS, exists := afterFiles[file]; exists {
			comparison := compareFileSummaries(beforeFS, afterFS)
			if comparison.NetChange < 0 {
				result.ImprovedFiles = append(result.ImprovedFiles, comparison)
			} else if comparison.NetChange > 0 {
				result.WorsenedFiles = append(result.WorsenedFiles, comparison)
			}
		} else {
			result.RemovedFiles = append(result.RemovedFiles, file)
		}
	}

	// Find new files
	for file := range afterFiles {
		if _, exists := beforeFiles[file]; !exists {
			result.NewFiles = append(result.NewFiles, file)
		}
	}

	return result
}

// compareFileSummaries compares two FileSummaries and returns a FileComparison
func compareFileSummaries(before, after FileSummary) FileComparison {
	comparison := FileComparison{
		File:          before.File,
		ImprovedRules: make([]RuleComparison, 0),
		WorsenedRules: make([]RuleComparison, 0),
		NewRules:      make([]string, 0),
		RemovedRules:  make([]string, 0),
	}

	// Create maps for easier lookup
	beforeRules := make(map[string]RuleSummary)
	afterRules := make(map[string]RuleSummary)

	for _, rs := range before.RuleSummaries {
		beforeRules[rs.RuleID] = rs
		comparison.TotalBefore += rs.Count
	}
	for _, rs := range after.RuleSummaries {
		afterRules[rs.RuleID] = rs
		comparison.TotalAfter += rs.Count
	}

	// Compare rules that exist in both summaries
	for ruleID, beforeRS := range beforeRules {
		if afterRS, exists := afterRules[ruleID]; exists {
			change := afterRS.Count - beforeRS.Count
			// Determine the severity by taking the higher severity
			severity := beforeRS.Severity
			if string(afterRS.Severity) == string(SeverityError) {
				severity = SeverityError
			}
			ruleComparison := RuleComparison{
				RuleID:      ruleID,
				CountBefore: beforeRS.Count,
				CountAfter:  afterRS.Count,
				Change:      change,
				Severity:    severity,
			}
			if change < 0 {
				comparison.ImprovedRules = append(comparison.ImprovedRules, ruleComparison)
			} else if change > 0 {
				comparison.WorsenedRules = append(comparison.WorsenedRules, ruleComparison)
			}
		} else {
			comparison.RemovedRules = append(comparison.RemovedRules, ruleID)
		}
	}

	// Find new rules
	for ruleID := range afterRules {
		if _, exists := beforeRules[ruleID]; !exists {
			comparison.NewRules = append(comparison.NewRules, ruleID)
		}
	}

	comparison.NetChange = comparison.TotalAfter - comparison.TotalBefore
	return comparison
}

// NewToolSummary creates a new ToolSummary from a slice of AnalysisResults
func NewToolSummary(results []AnalysisResult) *ToolSummary {
	if len(results) == 0 {
		return &ToolSummary{}
	}

	// Group by file
	fileMap := make(map[string][]AnalysisResult)
	for _, result := range results {
		fileMap[result.File] = append(fileMap[result.File], result)
	}

	summary := &ToolSummary{
		Tool:          results[0].Tool,
		FileSummaries: make([]FileSummary, 0, len(fileMap)),
	}

	// Create file summaries
	for file, fileResults := range fileMap {
		fileSummary := FileSummary{
			File:          file,
			RuleSummaries: make([]RuleSummary, 0),
		}

		// Group by rule ID
		ruleMap := make(map[string][]AnalysisResult)
		for _, result := range fileResults {
			ruleMap[result.RuleID] = append(ruleMap[result.RuleID], result)
		}

		// Create rule summaries
		for ruleID, ruleResults := range ruleMap {
			if len(ruleResults) == 0 {
				continue
			}

			ruleSummary := RuleSummary{
				RuleID:      ruleID,
				Description: ruleResults[0].Description,
				Severity:    ruleResults[0].Severity,
				Count:       len(ruleResults),
				Violations:  make([]Violation, 0, len(ruleResults)),
			}

			// Add violations
			for _, result := range ruleResults {
				ruleSummary.Violations = append(ruleSummary.Violations, Violation{
					Line:    result.Line,
					Column:  result.Column,
					Message: result.Message,
				})
			}

			fileSummary.RuleSummaries = append(fileSummary.RuleSummaries, ruleSummary)
		}

		summary.FileSummaries = append(summary.FileSummaries, fileSummary)
	}

	return summary
} 