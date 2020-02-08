package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/beatlabs/gomodctl/internal"
)

// ErrNoVersionAvailable is exported
var ErrNoVersionAvailable = errors.New("no version available")

// Checker is exported
type Checker struct {
	Ctx context.Context
}

// VulnerabilitiesCheck is exported
func (c *Checker) VulnerabilitiesCheck(path string, vulnerabilityCheck bool, jsonOutputCheck bool) (map[string]internal.VulnerabilityResult, error) {
	return getModAndVulnerabilitiesCheck(c.Ctx, path, vulnerabilityCheck, jsonOutputCheck)
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

// CheckForVulnerabilities function check for possible vulnerabilities.
func CheckForVulnerabilities(ctx context.Context, packages []packageResult) map[string]internal.VulnerabilityResult {

	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		result map[string]internal.VulnerabilityResult
	)

	doneCh := make(chan bool, 1)
	wg.Add(len(packages))
	result = make(map[string]internal.VulnerabilityResult, 0)
	for i := 0; i < len(packages); i++ {
		go func(i int) {
			defer wg.Done()
			goSecDir := packages[i].dir + "/./..."
			arg := []string{"-quiet", "-fmt=json", goSecDir}
			cmd := exec.CommandContext(ctx, "gosec", arg...)
			out, _ := cmd.CombinedOutput()
			output := string(out)
			var vr internal.VulnerabilityResult
			err := json.Unmarshal([]byte(output), &vr)
			if err == nil {
				mu.Lock()
				result[packages[i].path] = vr
				mu.Unlock()
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(doneCh)
	}()
	select {
	case <-doneCh:
		return result
	}
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

func getModAndVulnerabilitiesCheck(ctx context.Context, path string, vulnerabilityCheck bool, jsonOutputCheck bool) (map[string]internal.VulnerabilityResult, error) {
	parser := versionParser{ctx: ctx}
	vs := make(map[string]internal.VulnerabilityResult)
	results, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}

	if vulnerabilityCheck {
		vs = CheckForVulnerabilities(ctx, results)
		if jsonOutputCheck {
			r, err := json.Marshal(vs)
			if err != nil {
				fmt.Println(err)
			}
			err = ioutil.WriteFile("output.json", r, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return vs, nil
}
