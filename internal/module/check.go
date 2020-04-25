package module

import (
	"context"
	"errors"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/beatlabs/gomodctl/internal"
)

// ErrNoVersionAvailable is exported
var ErrNoVersionAvailable = errors.New("no version available")

// Checker is exported
type Checker struct {
	Ctx context.Context
}

// Check is exported.
func (c *Checker) Check(path string) (map[string]internal.CheckResult, error) {
	return getModAndFilter(c.Ctx, path, getLatestVersion)
}

func getLatestVersion(_ *semver.Version, versions []*semver.Version) (*semver.Version, error) {
	if len(versions) == 0 {
		return nil, ErrNoVersionAvailable
	}

	sort.Sort(semver.Collection(versions))

	lastVersion := versions[len(versions)-1]

	return lastVersion, nil
}

func getModAndFilter(ctx context.Context, path string, filter func(*semver.Version, []*semver.Version) (*semver.Version, error)) (map[string]internal.CheckResult, error) {
	parser := ModParser{ctx: ctx}

	results, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}

	checkResults := make(map[string]internal.CheckResult)

	for _, result := range results {
		latestVersion, err := filter(result.LocalVersion, result.availableVersions)

		checkResult := internal.CheckResult{
			LocalVersion: result.LocalVersion,
		}

		if err != nil {
			checkResult.Error = err
		}

		if latestVersion != nil {
			checkResult.LatestVersion = latestVersion
		}

		checkResults[result.Path] = checkResult
	}

	return checkResults, nil
}
