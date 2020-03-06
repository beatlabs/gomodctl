package license

//go:generate go run license_embedder.go

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/beatlabs/gomodctl/internal"
	"github.com/beatlabs/gomodctl/internal/module"
	"github.com/go-resty/resty/v2"
	"github.com/google/licenseclassifier"
	"github.com/mholt/archiver/v3"
)

const invalidLicense = "Can't find license"
const licenseFilename = "LICENSE"

// Checker is exported
type Checker struct {
	classifier    *licenseclassifier.License
	restClient    *resty.Client
	ctx           context.Context
	versionParser *module.VersionParser
}

// NewChecker is exported
func NewChecker(ctx context.Context) (*Checker, error) {
	license, err := licenseclassifier.New(licenseclassifier.DefaultConfidenceThreshold, licenseclassifier.ArchiveBytes(licenseDB))
	if err != nil {
		return nil, err
	}

	return &Checker{
		classifier:    license,
		restClient:    resty.New(),
		ctx:           ctx,
		versionParser: module.NewVersionParser(ctx),
	}, nil
}

// Type is exportesd
func (f *Checker) Type(moduleName, version string) (string, error) {
	v, err := f.getVersion(moduleName, version)
	if err != nil {
		return "", err
	}

	localModulePath := createLocalModulePath(moduleName, v)
	if _, err := os.Stat(localModulePath); !os.IsNotExist(err) {
		return f.getTypeFromLocalFile(localModulePath)
	}

	response, err := f.restClient.R().
		SetContext(f.ctx).
		Get(createGoProxyURLForVersion(moduleName, v))
	if err != nil {
		return "", err
	}

	if !response.IsSuccess() {
		return "", errors.New(response.String())
	}

	tempZipName := fmt.Sprintf("%s-%s.*.zip", strings.ReplaceAll(moduleName, "/", ""), v.Original())

	tempFile, err := ioutil.TempFile("", tempZipName)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(response.Body())
	if err != nil {
		return "", err
	}

	err = tempFile.Close()
	if err != nil {
		return "", err
	}

	match := invalidLicense
	err = archiver.Walk(tempFile.Name(), func(file archiver.File) error {
		if strings.HasPrefix(file.Name(), licenseFilename) {
			b, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			match = f.classifier.NearestMatch(string(b)).Name
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return match, err
}

func (f *Checker) Types() (map[string]internal.LicenseResult, error) {
	parse, err := f.versionParser.Parse("")
	if err != nil {
		return nil, err
	}

	m := make(map[string]internal.LicenseResult)

	for _, result := range parse {
		licenseResult := internal.LicenseResult{
			LocalVersion: result.LocalVersion,
		}

		licenseType, err := f.getTypeFromLocalFile(result.Dir)
		if err != nil {
			licenseResult.Error = err
		} else {
			licenseResult.Type = licenseType
		}

		m[result.Path] = licenseResult
	}

	return m, nil
}

func (f *Checker) getTypeFromLocalFile(path string) (string, error) {
	match := invalidLicense

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return match, err
	}

	for _, info := range dir {
		if strings.HasPrefix(info.Name(), licenseFilename) {
			b, err := ioutil.ReadFile(filepath.Join(path, info.Name()))
			if err != nil {
				return match, err
			}

			match = f.classifier.NearestMatch(string(b)).Name

			break
		}
	}

	return match, err
}

func (f *Checker) getVersion(moduleName, version string) (*semver.Version, error) {
	if version == "" {
		v, err := f.getLatestVersion(moduleName)
		if err != nil {
			return nil, err
		}

		return v, nil
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (f *Checker) getLatestVersion(moduleName string) (*semver.Version, error) {
	resp := &response{}

	response, err := f.restClient.R().
		SetContext(f.ctx).
		SetHeader("Accept", "application/json").
		SetResult(resp).
		Get(createGoProxyURLForLatestVersion(moduleName))
	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, errors.New(response.String())
	}

	if resp.Version == "" {
		return nil, fmt.Errorf("no version available, %s", response.String())
	}

	return semver.NewVersion(resp.Version)
}

func createGoProxyURLForLatestVersion(moduleName string) string {
	return fmt.Sprintf("%s/%s/@latest", getGoProxy(), moduleName)
}

func createLocalModulePath(moduleName string, version *semver.Version) string {
	return fmt.Sprintf("%s/pkg/mod/%s@%s", os.Getenv("GOPATH"), encodeModuleName(moduleName), version.Original())
}

func createGoProxyURLForVersion(moduleName string, version *semver.Version) string {
	return fmt.Sprintf("%s/%s/@v/%s.zip", getGoProxy(), encodeModuleName(moduleName), version.Original())
}

func encodeModuleName(moduleName string) string {
	regExp := regexp.MustCompile(`[[:upper:]]`)

	return regExp.ReplaceAllStringFunc(moduleName, func(s string) string {
		return "!" + strings.ToLower(s)
	})
}

func getGoProxy() string {
	goProxyEnv := os.Getenv("GOPROXY")

	goProxies := strings.Split(goProxyEnv, ",")

	n := 0
	for _, x := range goProxies {
		if x != "direct" && strings.TrimSpace(x) != "" {
			goProxies[n] = x
			n++
		}
	}
	goProxies = goProxies[:n]

	if len(goProxies) == 0 {
		return "https://proxy.golang.org"
	}

	return goProxies[0]
}

type response struct {
	Version string    `json:"Version"`
	Time    time.Time `json:"Time"`
}
