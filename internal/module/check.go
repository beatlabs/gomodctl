package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
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

// Check is exported
func (c *Checker) VulnerabilitiesCheck(path string, vulnerabilityCheck bool) (map[string]internal.VulnerabilityResult, error) {
	return getModAndVulnerabilitiesCheck(c.Ctx, path, vulnerabilityCheck)
}

// Check is exported
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

func vulnerabilityCheckFn(ctx context.Context, packages []packageResult) map[string]internal.VulnerabilityResult {

	result := make(map[string]internal.VulnerabilityResult)

	for _, r := range packages {
		fmt.Println("Scan: " + r.dir)
		goSecDir := r.dir + "/./..."
		arg := []string{"-quiet", "-fmt=json", goSecDir}
		cmd := exec.CommandContext(ctx, "/Users/thomaschavakis/go/bin/gosec", arg...)
		out, _ := cmd.CombinedOutput()
		output := string(out)
		var vr internal.VulnerabilityResult
		err := json.Unmarshal([]byte(output), &vr)
		if err != nil {
			result[r.path] = vr
		}
		result[r.path] = vr
	}
	return result
}

func getModAndFilter(ctx context.Context, path string, filter func(*semver.Version, []*semver.Version) (*semver.Version, error)) (map[string]internal.CheckResult, error) {
	parser := versionParser{ctx: ctx}

	results, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}

	checkResults := make(map[string]internal.CheckResult)

	for _, result := range results {
		latestVersion, err := filter(result.localVersion, result.availableVersions)

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

func getModAndVulnerabilitiesCheck(ctx context.Context, path string, vulnerabilityCheck bool) (map[string]internal.VulnerabilityResult, error) {
	parser := versionParser{ctx: ctx}
	vs := make(map[string]internal.VulnerabilityResult)
	results, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}

	if vulnerabilityCheck {
		vs = vulnerabilityCheckFn(ctx, results)
	}
	return vs, nil
}
