package internal

// CheckResult is exported.
type CheckResult struct {
	Name          string
	LocalVersion  string
	LatestVersion string
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
