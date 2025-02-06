package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	reportFile = flag.String("report", "report.json", "path to the semgrep report")
)

// SemgrepMetadata represents metadata about a finding
type SemgrepMetadata struct {
	Category   string   `json:"category"`
	Confidence string   `json:"confidence"`
	Impact     string   `json:"impact"`
	Likelihood string   `json:"likelihood"`
	CWE        []string `json:"cwe"`
	OWASP      []string `json:"owasp"`
	Technology []string `json:"technology"`
}

// Position represents a position in a file
type Position struct {
	Line   int `json:"line"`
	Col    int `json:"col"`
	Offset int `json:"offset,omitempty"`
}

// CliMatchExtra contains additional information about a finding
type CliMatchExtra struct {
	Message     string          `json:"message"`
	Metadata    SemgrepMetadata `json:"metadata"`
	Severity    string          `json:"severity"`
	Fingerprint string          `json:"fingerprint"`
	Lines       string          `json:"lines"`
	IsIgnored   bool            `json:"is_ignored"`
}

// CliMatch represents a single finding
type CliMatch struct {
	CheckID string        `json:"check_id"`
	Path    string        `json:"path"`
	Start   Position      `json:"start"`
	End     Position      `json:"end"`
	Extra   CliMatchExtra `json:"extra"`
}

// ScannedAndSkipped represents scanned and skipped paths
type ScannedAndSkipped struct {
	Scanned []string `json:"scanned"`
}

// SemgrepReport represents the full JSON report
type SemgrepReport struct {
	Version                string            `json:"version"`
	Results                []CliMatch        `json:"results"`
	Errors                 []interface{}     `json:"errors"`
	Paths                  ScannedAndSkipped `json:"paths"`
	InterfileLanguagesUsed []string          `json:"interfile_languages_used"`
	SkippedRules           []interface{}     `json:"skipped_rules"`
}

// severityEmoji maps severity levels to emojis
var severityEmoji = map[string]string{
	"ERROR":   "🚨",
	"WARNING": "🚧",
	"INFO":    "ℹ️",
	"UNKNOWN": "🤷",
}

// severityOrder defines the order in which severities should be displayed
var severityOrder = []string{"ERROR", "WARNING", "INFO", "UNKNOWN"}

func createMarkdownComment(results []CliMatch) string {
	if len(results) == 0 {
		return "✅ No security issues found!"
	}

	var markdown []string
	markdown = append(markdown, "# 🔒 Security Scan Results\n")

	// Group results by severity
	severityGroups := make(map[string][]CliMatch)
	for _, result := range results {
		severity := result.Extra.Severity
		if severity == "" {
			severity = "UNKNOWN"
		}
		severityGroups[severity] = append(severityGroups[severity], result)
	}

	// Add summary
	markdown = append(markdown, "## Summary\n")
	for _, severity := range severityOrder {
		if findings, ok := severityGroups[severity]; ok {
			count := len(findings)
			emoji := severityEmoji[severity]
			plural := "s"
			if count == 1 {
				plural = ""
			}
			markdown = append(markdown, fmt.Sprintf("- %s **%s**: %d finding%s", emoji, severity, count, plural))
		}
	}
	markdown = append(markdown, "")

	// Add detailed findings in dropdown
	markdown = append(markdown, "<details>")
	markdown = append(markdown, "<summary>Click to see detailed findings</summary>\n")

	// Add detailed findings
	markdown = append(markdown, "## Detailed Findings\n")
	for _, severity := range severityOrder {
		findings, ok := severityGroups[severity]
		if !ok {
			continue
		}

		markdown = append(markdown, fmt.Sprintf("### %s Severity\n", severity))

		for _, result := range findings {
			markdown = append(markdown, fmt.Sprintf("#### %s\n", result.CheckID))
			markdown = append(markdown, fmt.Sprintf("**File:** `%s:%d`\n", result.Path, result.Start.Line))
			markdown = append(markdown, fmt.Sprintf("**Message:** %s\n", result.Extra.Message))

			if result.Extra.Lines != "" {
				markdown = append(markdown, "```")
				markdown = append(markdown, result.Extra.Lines)
				markdown = append(markdown, "```\n")
			}

			markdown = append(markdown, "**Additional Info:**")
			markdown = append(markdown, fmt.Sprintf("- Category: %s", result.Extra.Metadata.Category))
			markdown = append(markdown, fmt.Sprintf("- Confidence: %s", result.Extra.Metadata.Confidence))
			markdown = append(markdown, fmt.Sprintf("- Impact: %s", result.Extra.Metadata.Impact))
			markdown = append(markdown, fmt.Sprintf("- Likelihood: %s", result.Extra.Metadata.Likelihood))

			if len(result.Extra.Metadata.CWE) > 0 {
				markdown = append(markdown, fmt.Sprintf("- CWE: %s", strings.Join(result.Extra.Metadata.CWE, ", ")))
			}
			if len(result.Extra.Metadata.OWASP) > 0 {
				markdown = append(markdown, fmt.Sprintf("- OWASP: %s", strings.Join(result.Extra.Metadata.OWASP, ", ")))
			}

			markdown = append(markdown, "")
		}
	}

	markdown = append(markdown, "</details>")
	return strings.Join(markdown, "\n")
}

func main() {
	flag.Parse()

	// Read the JSON file
	data, err := os.ReadFile(*reportFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the JSON data
	var report SemgrepReport
	if err := json.Unmarshal(data, &report); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Generate and print the markdown comment
	markdown := createMarkdownComment(report.Results)
	fmt.Println(markdown)
}
