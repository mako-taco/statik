---
layout: default
title: Statik - Static Analysis Tool Output Parser
description: A powerful CLI tool that unifies and compares outputs from various static analysis tools
---

# Statik

<div class="hero">
  <h1>Unify Your Static Analysis Workflow</h1>
  <p class="lead">Statik is a powerful CLI tool that parses and compares outputs from various static analysis tools, providing a unified interface for better code quality management.</p>
</div>

## 🚀 Features

- **Unified Output Format**: Convert outputs from different static analysis tools into a consistent JSON format
- **Smart Comparisons**: Compare analysis results before and after changes to track improvements or regressions
- **Custom Rule Summaries**: Intelligent handling of specific rules (e.g., ESLint's max-lines rule)
- **Multiple Tool Support**: Currently supports:
  - TypeScript Compiler (tsc)
  - ESLint
- **Flexible Input**: Read from files or directly from stdin
- **Configurable Severity**: Option to ignore warnings when comparing results

## 📦 Installation

```bash
go install github.com/statik@latest
```

## 🛠️ Usage

### Parse Command

Convert static analysis tool output into a unified JSON format:

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

Track changes in your code quality over time:

```bash
# Compare two summaries
statik compare before.json after.json

# Compare summaries and ignore warnings
statik compare before.json after.json --ignore-warnings
```

### Example Output

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

## 🔍 Smart Rule Handling

Statik provides intelligent handling of specific rules. For example, with ESLint's `max-lines` rule:

- Instead of just counting violations, it shows the actual difference between current and maximum allowed lines
- Makes it easier to track how far you are from the limit
- Helps prioritize which files need the most attention

## 🎯 Use Cases

- **CI/CD Integration**: Track code quality metrics in your pipeline
- **Pre-commit Hooks**: Prevent regressions before they're committed
- **Code Review**: Quickly identify problematic changes
- **Team Metrics**: Track code quality improvements over time
- **Project Health**: Get a unified view of your codebase's static analysis status

## 🤝 Contributing

We welcome contributions! Whether it's:

- Adding support for new static analysis tools
- Implementing custom rule summaries
- Improving documentation
- Bug fixes and feature enhancements

Check out our [Contributing Guide](CONTRIBUTING.md) to get started.

## 📄 License

MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [GitHub Repository](https://github.com/statik)
- [Issue Tracker](https://github.com/statik/issues)
- [Documentation](https://github.com/statik/docs)

---

<div class="cta">
  <h2>Ready to improve your static analysis workflow?</h2>
  <p>Get started with Statik today!</p>
  <a href="https://github.com/statik" class="button">View on GitHub</a>
</div>

<style>
.hero {
  text-align: center;
  padding: 4rem 2rem;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 8px;
  margin-bottom: 2rem;
}

.hero h1 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  color: #2c3e50;
}

.lead {
  font-size: 1.25rem;
  color: #34495e;
}

.cta {
  text-align: center;
  padding: 3rem 2rem;
  background: #f8f9fa;
  border-radius: 8px;
  margin-top: 3rem;
}

.button {
  display: inline-block;
  padding: 0.8rem 1.5rem;
  background: #2c3e50;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  margin-top: 1rem;
  transition: background 0.3s ease;
}

.button:hover {
  background: #34495e;
}
</style>
