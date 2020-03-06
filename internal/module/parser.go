package module

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
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

func NewVersionParser(ctx context.Context) *VersionParser {
	return &VersionParser{ctx: ctx}
}

type VersionParser struct {
	ctx context.Context
}

type PackageResult struct {
	Path              string
	LocalVersion      *semver.Version
	availableVersions []*semver.Version
	Dir               string
}

// Parse is exported
func (v *VersionParser) Parse(path string) ([]PackageResult, error) {
	cmd := exec.CommandContext(v.ctx, "go", "list", "-m", "-versions", "-json", "all")

	if path != "" {
		home := viper.GetString("home")

		if strings.HasPrefix(path, home) {
			l := path[len(home):]
			cmd.Dir = filepath.Join(home, l)
		} else {
			cmd.Dir = filepath.Join(home, path)
		}
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

	var result []PackageResult

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

			result = append(result, PackageResult{
				Path:              it.Path,
				LocalVersion:      semver.MustParse(it.Version),
				Dir:               it.Dir,
				availableVersions: availableVersions,
			})
		}
	}

	return result, nil
}
