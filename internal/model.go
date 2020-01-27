package internal

import "github.com/Masterminds/semver"

// CheckResult is exported.
type CheckResult struct {
	LocalVersion  *semver.Version
	LatestVersion *semver.Version
	Error         error
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
