package module

import (
	"strings"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/prometheus/common/log"
	"github.com/tcnksm/go-latest"
)

// Checker is exported.
type Checker struct {
}

// Check is exported.
func (c *Checker) Check(path string) ([]internal.CheckResult, error) {
	parser := VersionParser{GoModPath: "go.mod"}

	results, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}

	var checkResults []internal.CheckResult

	for _, result := range results {
		fetcher := factory(result.Path, result.LocalVersion)

		latestVersion, err := fetcher.Fetch()
		if err != nil {
			log.Warn(err)
			latestVersion = "unknown"
		}

		checkResults = append(checkResults, internal.CheckResult{
			Name:          result.Path,
			LocalVersion:  result.LocalVersion,
			LatestVersion: latestVersion,
		})
	}

	return checkResults, nil
}

// VersionFetcher is exported.
type VersionFetcher interface {
	Fetch() (string, error)
}

func factory(path string, version string) VersionFetcher {
	if strings.HasPrefix(path, "github.com/") {
		return &githubFetcher{
			path:         path,
			localVersion: version,
		}
	}

	return &DummyFetcher{
		Path: path,
	}
}

// DummyFetcher is exported.
type DummyFetcher struct {
	Path string
}

// Fetch is exported.
func (f *DummyFetcher) Fetch() (string, error) {
	return "unknown", nil
}

type githubFetcher struct {
	path         string
	localVersion string
}

// Fetch is exported.
func (g *githubFetcher) Fetch() (string, error) {
	usernameAndRepo := strings.TrimPrefix(g.path, "github.com/")

	splits := strings.Split(usernameAndRepo, "/")

	githubTag := &latest.GithubTag{
		Owner:      splits[0],
		Repository: splits[1],
	}

	check, err := latest.Check(githubTag, g.localVersion)
	if err != nil {
		return "", err
	}

	if check.Outdated {
		return check.Current, nil
	}

	return g.localVersion, nil
}
