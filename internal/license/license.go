package license

import (
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
}

// NewChecker is exported
func NewChecker() (*Checker, error) {
	license, err := licenseclassifier.New(licenseclassifier.DefaultConfidenceThreshold)
	if err != nil {
		return nil, err
	}

	return &Checker{
		classifier: license,
		restClient: resty.New(),
	}, nil
}

// Type is exportesd
func (f *Checker) Type(moduleName, version string) (string, error) {
	v, err := f.getVersion(moduleName, version)
	if err != nil {
		return "", err
	}

	response, err := f.restClient.R().
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
		SetHeader("Accept", "application/json").
		SetResult(resp).
		Get(createGoProxyURLForLatestVersion(moduleName))
	if err != nil {
		return nil, err
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
	goProxy := os.Getenv("GOPROXY")
	if len(goProxy) == 0 || goProxy == "direct" {
		return "https://proxy.golang.org"
	}

	return goProxy
}

type response struct {
	Version string    `json:"Version"`
	Time    time.Time `json:"Time"`
}
