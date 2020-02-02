package module

import (
	"io/ioutil"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/beatlabs/gomodctl/internal"
	"golang.org/x/mod/modfile"
)

type result struct {
	localVersion      *semver.Version
	availableVersions []*semver.Version
	err               error
}

func modFileToError(modFile *modfile.File) (map[string]result, error) {
	checkResults := make(map[string]result)

	for _, require := range modFile.Require {
		if !require.Indirect {
			checkResult := result{}

			version, err := semver.NewVersion(require.Mod.Version)
			if err != nil {
				checkResult.err = err
			} else {
				checkResult.localVersion = version
			}

			checkResults[require.Mod.Path] = checkResult
		}
	}

	return checkResults, nil
}

func parseMod(path string) (*modfile.File, []byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	parse, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, nil, err
	}

	return parse, content, nil
}

func checkForVersions(versions map[string]result) (map[string]result, error) {
	return versions, nil
}

func GetLatest(path string) (map[string]internal.CheckResult, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	absoluteFile := filepath.Join(absolutePath, goMod)

	modFile, _, err := parseMod(absoluteFile)
	if err != nil {
		return nil, err
	}

	filter := getLatestVersion

	versions, err := modFileToError(modFile)
	if err != nil {
		return nil, err
	}

	versions, err = checkForVersions(versions)
	if err != nil {
		return nil, err
	}

	return findLastVersion(versions, filter), nil
}

func UpdateMod(path string) (map[string]internal.CheckResult, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	absoluteFile := filepath.Join(absolutePath, goMod)
	backupFile := filepath.Join(absolutePath, goModBackup)

	modFile, originalContent, err := parseMod(absoluteFile)
	if err != nil {
		return nil, err
	}

	filter := getLatestVersion

	versions, err := modFileToError(modFile)
	if err != nil {
		return nil, err
	}

	versions, err = checkForVersions(versions)
	if err != nil {
		return nil, err
	}

	latestMinors := findLastVersion(versions, filter)

	updates := 0

	for moduleName, result := range latestMinors {
		if result.Error == nil && result.LatestVersion.GreaterThan(result.LocalVersion) {
			err := modFile.DropRequire(moduleName)
			if err != nil {
				return nil, err
			}

			err = modFile.AddRequire(moduleName, result.LatestVersion.Original())
			if err != nil {
				return nil, err
			}

			updates++
		}
	}

	if updates > 0 {
		format, err := modFile.Format()
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(absoluteFile, format, 0666)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(backupFile, originalContent, 0666)
		if err != nil {
			return nil, err
		}
	}

	return latestMinors, nil
}

func findLastVersion(versions map[string]result, filter func(*semver.Version, []*semver.Version) (*semver.Version, error)) map[string]internal.CheckResult {
	checkResults := make(map[string]internal.CheckResult)

	for path, result := range versions {
		checkResult := internal.CheckResult{
			LocalVersion: result.localVersion,
		}
		if result.err != nil {
			checkResult.Error = result.err
		} else {
			latestVersion, err := filter(result.localVersion, result.availableVersions)
			if err != nil {
				checkResult.Error = err
			}
			if latestVersion != nil {
				checkResult.LatestVersion = latestVersion
			}
		}

		checkResults[path] = checkResult
	}

	return checkResults
}
