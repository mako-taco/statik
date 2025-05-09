---
layout: default
title: Statik - Static Analysis Tool Output Parser
description: A powerful CLI tool that helps teams gradually introduce new static analysis rules without blocking development
---

# Statik

<div class="hero">
  <h1>Introduce New Static Analysis Rules Safely</h1>
  <p class="lead">Statik helps teams gradually adopt new static analysis rules by tracking progress over time, without requiring immediate compliance across the entire codebase.</p>
</div>

## 🎯 The Problem

When introducing new static analysis rules to a large codebase, teams often face a difficult choice:

- **Option 1**: Fix all violations immediately (blocks development, high risk)
- **Option 2**: Disable the rule (misses potential issues)
- **Option 3**: Use Statik (track progress, ensure improvement)

## 💡 The Solution

Statik provides a third way: gradually improve code quality while ensuring no regressions. It helps you:

- Track the number of violations over time
- Ensure new code doesn't introduce more violations
- Set and monitor progress towards zero violations
- Compare before/after changes to catch regressions

## 🚀 Key Features

- **Progress Tracking**: Monitor the number of violations over time
- **Regression Prevention**: Fail CI if new violations are introduced
- **Smart Comparisons**: Compare analysis results before and after changes
- **Custom Rule Summaries**: Intelligent handling of specific rules (e.g., ESLint's max-lines rule)
- **Multiple Tool Support**: Currently supports:
  - TypeScript Compiler (tsc)
  - ESLint

## 🔧 Supported Tools

| Tool                      | Description                                     | File Extensions                      | Custom Rule Support          |
| ------------------------- | ----------------------------------------------- | ------------------------------------ | ---------------------------- |
| TypeScript Compiler (tsc) | TypeScript's built-in type checker and compiler | `.ts`, `.tsx`                        | No custom rules yet          |
| ESLint                    | JavaScript/TypeScript linter                    | `.js`, `.jsx`, `.ts`, `.tsx`, `.vue` | Yes (e.g., `max-lines` rule) |

_More tools coming soon! [Request a tool](https://github.com/statik/issues) or [contribute](CONTRIBUTING.md) your own parser._

## 📦 Installation

```bash
go install github.com/statik@latest
```

## 🛠️ Usage

### 1. Initial Setup

```bash
# Generate initial baseline
eslint --format json . | statik parse eslint > baseline.json
```

### 2. Track Progress

```bash
# Generate current state
eslint --format json . | statik parse eslint > current.json

# Compare with baseline
statik compare baseline.json current.json
```

### 3. CI Integration

```bash
# Fail if new violations are introduced
statik compare baseline.json current.json

# Or only fail on error-level regressions
statik compare baseline.json current.json --ignore-warnings
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
      "total_before": 2,
      "total_after": 0,
      "net_change": -2
    }
  ],
  "worsened_files": [
    {
      "file": "src/file2.ts",
      "worsened_rules": [
        {
          "rule_id": "TS2322",
          "count_before": 0,
          "count_after": 1,
          "change": 1
        }
      ],
      "total_before": 0,
      "total_after": 1,
      "net_change": 1
    }
  ]
}
```

## 🔍 Smart Rule Handling

Statik provides intelligent handling of specific rules. For example, with ESLint's `max-lines` rule:

- Shows the actual difference between current and maximum allowed lines
- Makes it easier to track progress towards the limit
- Helps prioritize which files need the most attention

## 🎯 Use Cases

- **Gradual Rule Adoption**: Introduce new rules without blocking development
- **CI/CD Integration**: Ensure no new violations are introduced
- **Team Metrics**: Track progress towards zero violations
- **Code Review**: Quickly identify problematic changes
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
  <h2>Ready to improve your code quality gradually?</h2>
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
