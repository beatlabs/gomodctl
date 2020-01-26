package module

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/beatlabs/gomodctl/internal"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// Updater is exported
type Updater struct {
	Ctx context.Context
}

const goMod = "go.mod"

// Update is exported
func (u *Updater) Update(path string, updateType internal.UpdateType) (map[string]internal.CheckResult, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	absoluteFile := filepath.Join(absolutePath, goMod)

	content, err := ioutil.ReadFile(absoluteFile)
	if err != nil {
		return nil, err
	}

	parse, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, err
	}

	filter := getLatestVersion
	if updateType == internal.Minor {
		filter = getLatestMinorVersion
	}

	latestMinors, err := getModAndFilter(u.Ctx, absolutePath, filter)
	if err != nil {
		return nil, err
	}

	var updates []*modfile.Require

	for moduleName, result := range latestMinors {
		if result.Error == nil && result.LatestVersion.GreaterThan(result.LocalVersion) {
			updates = append(updates, &modfile.Require{
				Mod: module.Version{
					Path:    moduleName,
					Version: result.LatestVersion.Original(),
				},
				Indirect: false,
			})
		}
	}

	if len(updates) > 0 {
		parse.SetRequire(updates)

		format, err := parse.Format()
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(absoluteFile, format, 0666)
		if err != nil {
			return nil, err
		}
	}

	return latestMinors, nil
}

func getLatestMinorVersion(current *semver.Version, versions []*semver.Version) (*semver.Version, error) {
	n := 0
	for _, version := range versions {
		if version.Major() == current.Major() {
			versions[n] = version
			n++
		}
	}
	versions = versions[:n]

	return getLatestVersion(nil, versions)
}
