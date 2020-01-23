package module

import (
	"errors"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/beatlabs/gomodctl/internal"
)

// ErrNoVersionAvailable is exported
var ErrNoVersionAvailable = errors.New("no version available")

// Checker is exported
type Checker struct {
}

// Check is exported
func (c *Checker) Check(path string) (map[string]internal.CheckResult, error) {
	parser := versionParser{}

	results, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}

	checkResults := make(map[string]internal.CheckResult)

	for _, result := range results {
		latestVersion, err := getLatestVersion(result.availableVersions)

		checkResult := internal.CheckResult{
			LocalVersion: result.localVersion,
		}

		if err != nil {
			checkResult.Error = err
		}

		if latestVersion != nil {
			checkResult.LatestVersion = latestVersion
		}

		checkResults[result.path] = checkResult
	}

	return checkResults, nil
}

func getLatestVersion(versions []*semver.Version) (*semver.Version, error) {
	if len(versions) == 0 {
		return nil, ErrNoVersionAvailable
	}

	sort.Sort(semver.Collection(versions))

	lastVersion := versions[len(versions)-1]

	return lastVersion, nil
}
