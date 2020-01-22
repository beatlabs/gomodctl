package module

import (
	"strings"

	"github.com/prometheus/common/log"
	"github.com/beatlabs/gomodctl/internal"
	"github.com/tcnksm/go-latest"
)

type Checker struct {
}

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

type VersionFetcher interface {
	Fetch() (string, error)
}

func factory(path string, version string) VersionFetcher {
	if strings.HasPrefix(path, "github.com/") {
		return &githubFetcher{
			path:         path,
			localVersion: version,
		}
	} else {
		return &DummyFetcher{
			Path: path,
		}
	}
}

type DummyFetcher struct {
	Path string
}

func (f *DummyFetcher) Fetch() (string, error) {
	return "unknown", nil
}

type githubFetcher struct {
	path         string
	localVersion string
}

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
