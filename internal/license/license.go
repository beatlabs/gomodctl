package license

//go:generate go run license_embedder.go

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/go-resty/resty/v2"
	"github.com/google/licenseclassifier"
	"github.com/mholt/archiver/v3"
)

const invalidLicense = "Can't find license"

// Checker is exported
type Checker struct {
	classifier *licenseclassifier.License
	restClient *resty.Client
	ctx        context.Context
}

// NewChecker is exported
func NewChecker(ctx context.Context) (*Checker, error) {
	license, err := licenseclassifier.New(licenseclassifier.DefaultConfidenceThreshold, licenseclassifier.ArchiveBytes(licenseDB))
	if err != nil {
		return nil, err
	}

	return &Checker{
		classifier: license,
		restClient: resty.New(),
		ctx:        ctx,
	}, nil
}

// Type is exportesd
func (f *Checker) Type(moduleName, version string) (string, error) {
	v, err := f.getVersion(moduleName, version)
	if err != nil {
		return "", err
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
		if strings.HasPrefix(file.Name(), "LICENSE") {
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
