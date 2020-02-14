package module

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/spf13/viper"
)

var regex = regexp.MustCompile(`({([^}]*)})`)

type item struct {
	Path     string   `json:"Path"`
	Version  string   `json:"Version"`
	Versions []string `json:"Versions"`
	Indirect bool     `json:"Indirect"`
	Main     bool     `json:"Main"`
	Dir      string   `json:"Dir"`
	GoMod    string   `json:"GoMod"`
}

type versionParser struct {
	ctx context.Context
}

type packageResult struct {
	path              string
	localVersion      *semver.Version
	availableVersions []*semver.Version
	dir               string
}

// Parse is exported
func (v *versionParser) Parse(path string) ([]packageResult, error) {
	cmd := exec.CommandContext(v.ctx, "go", "list", "-m", "-versions", "-json", "all")

	if path != "" {
		home := viper.GetString("home")
		cmd.Dir = filepath.Join(home, path)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return nil, fmt.Errorf("with output [%s] %w", out, err)
		}

		return nil, err
	}

	output := string(out)
	versionOutputs := regex.FindAllString(output, -1)

	var result []packageResult

	for _, versionOutput := range versionOutputs {
		it := item{}

		err := json.Unmarshal([]byte(versionOutput), &it)
		if err != nil {
			return nil, err
		}

		if !it.Indirect && !it.Main {
			availableVersions := make([]*semver.Version, len(it.Versions))

			for i, version := range it.Versions {
				availableVersions[i] = semver.MustParse(version)
			}

			result = append(result, packageResult{
				path:              it.Path,
				localVersion:      semver.MustParse(it.Version),
				dir:               it.Dir,
				availableVersions: availableVersions,
			})
		}
	}

	return result, nil
}
