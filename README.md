# Statik

A CLI tool for parsing and comparing static analysis tool outputs.

## Features

- Parse outputs from various static analysis tools into a unified JSON format
- Compare static analysis results before and after changes
- Support for TypeScript compiler and ESLint outputs
- Read from files or stdin
- Exit with code 1 if any files have worsened (configurable)

## Installation

```bash
go install github.com/statik@latest
```

## Usage

### Parse Command

Parse static analysis tool output into a unified JSON format:

```bash
# Parse TypeScript compiler output from a file
statik parse tsc tsc-output.txt > summary.json

# Parse TypeScript compiler output directly from stdin
tsc --noEmit | statik parse tsc > summary.json

# Parse ESLint output from a file
statik parse eslint eslint-output.json > summary.json

# Parse ESLint output directly from stdin
eslint --format json . | statik parse eslint > summary.json
```

### Compare Command

Compare two static analysis summaries to see what has improved or worsened:

```bash
# Compare two summaries
statik compare before.json after.json

# Compare summaries and ignore warnings
statik compare before.json after.json --ignore-warnings
```

The compare command will:

1. Output a JSON comparison showing:
   - Files that improved or worsened
   - New files added
   - Files removed
   - Per-file breakdown of rule changes
2. Exit with code 1 if any files have worsened, unless:
   - The `--ignore-warnings` flag is set, in which case it only exits with code 1 if there are error-level regressions

Example output:

```json
{
  "improved_files": [
    {
      "file": "src/file1.ts",
      "improved_rules": [
        {
          "rule_id": "TS2345",
          "count_before": 2,
          "count_after": 0,
          "change": -2
        }
      ],
      "worsened_rules": [],
      "new_rules": [],
      "removed_rules": [],
      "total_before": 2,
      "total_after": 0,
      "net_change": -2
    }
  ],
  "worsened_files": [
    {
      "file": "src/file2.ts",
      "improved_rules": [],
      "worsened_rules": [
        {
          "rule_id": "TS2322",
          "count_before": 0,
          "count_after": 1,
          "change": 1
        }
      ],
      "new_rules": [],
      "removed_rules": [],
      "total_before": 0,
      "total_after": 1,
      "net_change": 1
    }
  ],
  "new_files": ["src/file3.ts"],
  "removed_files": ["src/file4.ts"]
}
```

### List Command

List available parsers:

```bash
statik list
```

## Supported Tools

- TypeScript Compiler (tsc)
- ESLint

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT
