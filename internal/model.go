package internal

import (
	"github.com/Masterminds/semver"
)

// CheckResult is exported.
type CheckResult struct {
	LocalVersion  *semver.Version
	LatestVersion *semver.Version
	Error         error
}

// LicenseResult is result for license check.
type LicenseResult struct {
	LocalVersion *semver.Version
	Type         string
	Error        error
}

// SearchResult is exported.
type SearchResult struct {
	Name        string
	Path        string
	ImportCount int
	Stars       int
	Score       float64
	Synopsis    string
}

// VulnerabilityResult is exported
type VulnerabilityResult struct {
	Issues []struct {
		Code       string `json:"code"`
		File       string `json:"file"`
		Line       string `json:"line"`
		Column     string `json:"column"`
		Details    string `json:"details"`
		RuleID     string `json:"rule_id"`
		Severity   string `json:"severity"`
		Confidence string `json:"confidence"`
		Cwe        struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		} `json:"cwe"`
	} `json:"issues"`
	Stats struct {
		Files int `json:"files"`
		Found int `json:"found"`
		Lines int `json:"lines"`
	} `json:"stats"`
}
