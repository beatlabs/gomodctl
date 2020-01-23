package internal

// CheckResult is exported.
type CheckResult struct {
	LocalVersion  Version
	LatestVersion Version
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

// Version is exported.
type Version interface {
	Original() string
	Major() int64
	Minor() int64
	Patch() int64
	Prerelease() string
	Metadata() string
}
