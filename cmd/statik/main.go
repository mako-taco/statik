package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/statik/pkg/parsers/checkstyle"
	"github.com/statik/pkg/parsers/eslint"
	"github.com/statik/pkg/parsers/tsc"
	"github.com/statik/pkg/plugin"
)

var (
	registry = plugin.NewRegistry()
	rootCmd  = &cobra.Command{
		Use:   "statik",
		Short: "Static analysis tool output parser",
		Long:  `A CLI tool that parses outputs from various static analysis tools and provides a unified interface.`,
	}

	parseCmd = &cobra.Command{
		Use:   "parse [parser-name] [input-file]",
		Short: "Parse static analysis tool output",
		Long: `Parse static analysis tool output. If no input file is provided, reads from stdin.
Example:
  # Parse from a file
  statik parse tsc output.txt

  # Parse from stdin
  tsc --noEmit | statik parse tsc`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires at least a parser name")
			}
			if len(args) > 2 {
				return fmt.Errorf("accepts at most 2 args, received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			parserName := args[0]
			var inputFile *os.File
			var err error

			parser, err := registry.GetParser(parserName)
			if err != nil {
				return fmt.Errorf("failed to get parser: %w", err)
			}

			// If no file path provided, read from stdin
			if len(args) == 1 {
				inputFile = os.Stdin
			} else {
				inputFile, err = os.Open(args[1])
				if err != nil {
					return fmt.Errorf("failed to open input file: %w", err)
				}
				defer inputFile.Close()
			}

			results, err := parser.Parse(inputFile)
			if err != nil {
				return fmt.Errorf("failed to parse input: %w", err)
			}

			// Create summary from results
			summary := plugin.NewToolSummary(results)

			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(summary)
		},
	}

	compareCmd = &cobra.Command{
		Use:   "compare [before-summary.json] [after-summary.json]",
		Short: "Compare two static analysis summaries",
		Long: `Compare two static analysis summaries to see what has improved or worsened.
The command takes two JSON summary files as input and outputs a comparison in JSON format.
If any files have worsened, the command will exit with code 1 unless --ignore-warnings is set.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the before summary
			beforeFile, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("failed to open before summary file: %w", err)
			}
			defer beforeFile.Close()

			var beforeSummary plugin.ToolSummary
			if err := json.NewDecoder(beforeFile).Decode(&beforeSummary); err != nil {
				return fmt.Errorf("failed to decode before summary: %w", err)
			}

			// Read the after summary
			afterFile, err := os.Open(args[1])
			if err != nil {
				return fmt.Errorf("failed to open after summary file: %w", err)
			}
			defer afterFile.Close()

			var afterSummary plugin.ToolSummary
			if err := json.NewDecoder(afterFile).Decode(&afterSummary); err != nil {
				return fmt.Errorf("failed to decode after summary: %w", err)
			}

			// Compare the summaries
			comparison := beforeSummary.Compare(&afterSummary)
			if comparison == nil {
				return fmt.Errorf("cannot compare summaries from different tools")
			}

			// Output the comparison in JSON format
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(comparison); err != nil {
				return fmt.Errorf("failed to encode comparison: %w", err)
			}

			// Check if any files have worsened
			ignoreWarnings, _ := cmd.Flags().GetBool("ignore-warnings")
			if len(comparison.WorsenedFiles) > 0 {
				// If --ignore-warnings is set, only exit with code 1 if there are error-level regressions
				if ignoreWarnings {
					hasErrorRegressions := false
					for _, file := range comparison.WorsenedFiles {
						for _, rule := range file.WorsenedRules {
							if rule.Severity == plugin.SeverityError {
								hasErrorRegressions = true
								break
							}
						}
						if hasErrorRegressions {
							break
						}
					}
					if hasErrorRegressions {
						os.Exit(1)
					}
				} else {
					os.Exit(1)
				}
			}

			return nil
		},
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available parsers",
		Run: func(cmd *cobra.Command, args []string) {
			parsers := registry.ListParsers()
			for _, name := range parsers {
				fmt.Println(name)
			}
		},
	}
)

func init() {
	// Register available parsers
	registry.Register(&tsc.Parser{})
	registry.Register(&eslint.Parser{})
	registry.Register(&checkstyle.Parser{})

	// Add commands
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(compareCmd)
	rootCmd.AddCommand(listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
} 