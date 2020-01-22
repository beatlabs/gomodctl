package internal

type CheckResult struct {
	Name          string
	LocalVersion  string
	LatestVersion string
}

type SearchResult struct {
	Name        string
	Path        string
	ImportCount int
	Stars       int
	Score       float64
	Synopsis    string
}
