package module

import (
	"context"
	"encoding/json"
	"os/exec"
	"sync"

	"github.com/beatlabs/gomodctl/internal"
)

// Scanner is exported.
type Scanner struct {
	Ctx context.Context
}

// Scan is exported.
func (c *Scanner) Scan(path string) (map[string]internal.VulnerabilityResult, error) {
	return getModAndVulnerabilitiesCheck(c.Ctx, path)
}

func getModAndVulnerabilitiesCheck(ctx context.Context, path string) (map[string]internal.VulnerabilityResult, error) {
	parser := VersionParser{ctx: ctx}
	vs := make(map[string]internal.VulnerabilityResult)
	results, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}

	vs = vulnerabilityScan(ctx, results)
	return vs, nil
}

// vulnerabilityScan function check for possible vulnerabilities using the gosec tool
func vulnerabilityScan(ctx context.Context, packages []PackageResult) map[string]internal.VulnerabilityResult {

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
			goSecDir := packages[i].Dir + "/./..."
			arg := []string{"-quiet", "-fmt=json", goSecDir}
			cmd := exec.CommandContext(ctx, "gosec", arg...)
			out, _ := cmd.CombinedOutput()
			output := string(out)
			var vr internal.VulnerabilityResult
			err := json.Unmarshal([]byte(output), &vr)
			if err == nil {
				mu.Lock()
				result[packages[i].Path] = vr
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
