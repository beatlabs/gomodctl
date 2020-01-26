package module

import (
	"context"
	"encoding/json"
	"os/exec"
	"regexp"
	"strings"

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

	home := viper.GetString("home")
	if path != "" {
		if home != "" && strings.Contains(path, "~") {
			path = home + strings.Replace(path, "~", "", 1) + "/"
		} else {
			path = path + "/"
		}

		cmd.Dir = path
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
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
