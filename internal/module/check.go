package module

import (
	"context"
	"errors"
	"sort"

	"github.com/spf13/viper"
	"github.com/Masterminds/semver"
	"github.com/beatlabs/gomodctl/internal"
)

// ErrNoVersionAvailable is exported
var ErrNoVersionAvailable = errors.New("no version available")
var ErrModuleIgnored = errors.New("module ignored")

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
	parser := VersionParser{ctx: ctx}

	results, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}

	ignoredModules := getIgnoredModules()

	checkResults := make(map[string]internal.CheckResult)

	for _, result := range results {
		checkResult := internal.CheckResult{
			LocalVersion: result.LocalVersion,
		}

		_, isIgnored := ignoredModules[result.path]
		if isIgnored {
			checkResult.Error = ErrModuleIgnored
		} else {
			latestVersion, err := filter(result.localVersion, result.availableVersions)

			if err != nil {
				checkResult.Error = err
			}

			if latestVersion != nil {
				checkResult.LatestVersion = latestVersion
			}
		}

		checkResults[result.Path] = checkResult
	}

	return checkResults, nil
}

type void struct{}

var member void

func getIgnoredModules() map[string]void {
	s := make(map[string]void)

	im := viper.GetStringSlice("ignored_modules")
	for _, m := range im {
		s[m] = member
	}

	return s
}
