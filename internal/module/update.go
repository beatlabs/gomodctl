package module

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/beatlabs/gomodctl/internal"
	"golang.org/x/mod/modfile"
)

// Updater is exported
type Updater struct {
	Ctx context.Context
}

const (
	goMod       = "go.mod"
	goModBackup = "go.mod.backup"
)

// Update is exported
func (u *Updater) Update(path string) (map[string]internal.CheckResult, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	absoluteFile := filepath.Join(absolutePath, goMod)
	backupFile := filepath.Join(absolutePath, goModBackup)

	content, err := ioutil.ReadFile(absoluteFile)
	if err != nil {
		return nil, err
	}

	parse, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, err
	}

	filter := getLatestVersion

	latestMinors, err := getModAndFilter(u.Ctx, absolutePath, filter)
	if err != nil {
		return nil, err
	}

	updates := 0

	for moduleName, result := range latestMinors {
		if result.Error == nil && result.LatestVersion.GreaterThan(result.LocalVersion) {
			err := parse.DropRequire(moduleName)
			if err != nil {
				return nil, err
			}

			err = parse.AddRequire(moduleName, result.LatestVersion.Original())
			if err != nil {
				return nil, err
			}

			updates++
		}
	}

	if updates > 0 {
		format, err := parse.Format()
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(absoluteFile, format, 0666)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(backupFile, content, 0666)
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
