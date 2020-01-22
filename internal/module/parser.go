package module

import (
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/mod/modfile"
)

// VersionParser is exported.
type VersionParser struct {
	GoModPath string
}

// PackageResult is exported.
type PackageResult struct {
	Path         string
	LocalVersion string
}

// Parse is exported.
func (v *VersionParser) Parse(path string) ([]PackageResult, error) {
	home := viper.GetString("home")
	if path != "" {
		if home != "" && strings.Contains(path, "~") {
			path = home + strings.Replace(path, "~", "", 1) + "/" + v.GoModPath
		} else {
			path = path + "/" + v.GoModPath
		}
	} else {
		path = v.GoModPath
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parse, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, err
	}

	var result []PackageResult

	for _, require := range parse.Require {
		if !require.Indirect {
			result = append(result, PackageResult{
				Path:         require.Mod.Path,
				LocalVersion: require.Mod.Version,
			})
		}
	}

	return result, nil
}
