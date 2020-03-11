package license

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

var content = []byte(`module github.com/beatlabs/gomodctl

go 1.13

require (
	github.com/Masterminds/semver v1.5.0
	github.com/frankban/quicktest v1.7.2 // indirect
	github.com/go-resty/resty/v2 v2.1.0
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/google/licenseclassifier v0.0.0-20200108231022-9dfa8d8474eb
	github.com/klauspost/compress v1.9.8 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mholt/archiver/v3 v3.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/pierrec/lz4 v2.4.1+incompatible // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/mod v0.2.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)`)

type LicenseTestSuite struct {
	suite.Suite
	tempFile string
	tempDir  string
}

func TestLicenseTestSuite(t *testing.T) {
	suite.Run(t, new(LicenseTestSuite))
}

func (s *LicenseTestSuite) SetupTest() {
	tmpddir, err := ioutil.TempDir("", "test")
	s.NoError(err)
	s.tempDir = tmpddir

	s.tempFile = filepath.Join(s.tempDir, "go.mod")

	err = ioutil.WriteFile(s.tempFile, content, 0666)
	s.NoError(err)
}

func (s *LicenseTestSuite) TearDownTest() {
	os.Remove(s.tempFile)
	os.RemoveAll(s.tempDir)
}

func (s *LicenseTestSuite) Test_GetLatestVersion() {
	checker, _ := NewChecker(context.TODO())

	latestVersion, err := checker.getLatestVersion("github.com/beatlabs/patron")

	s.NoError(err)
	s.NotNil(latestVersion)
	s.GreaterOrEqual(latestVersion.Major(), int64(0))
	s.GreaterOrEqual(latestVersion.Minor(), int64(30))
	s.GreaterOrEqual(latestVersion.Patch(), int64(0))
}

func (s *LicenseTestSuite) Test_GetLatestVersionOfNonExistingPackage() {
	checker, _ := NewChecker(context.TODO())

	_, err := checker.getLatestVersion("github.com/beatlabs/patro")

	s.Error(err)
}

func (s *LicenseTestSuite) Test_GetLicenseOfLatest() {
	checker, _ := NewChecker(context.TODO())

	latestVersion, err := checker.Type("github.com/beatlabs/patron", "")

	s.NoError(err)
	s.Equal("Apache-2.0", latestVersion)
}

func (s *LicenseTestSuite) Test_GetLicenseOfNonExistingVersion() {
	checker, _ := NewChecker(context.TODO())

	latestVersion, err := checker.Type("github.com/beatlabs/patron", "v999.0.0")

	s.EqualError(err, "not found: github.com/beatlabs/patron@v999.0.0: invalid version: unknown revision v999.0.0")
	s.Empty(latestVersion)
}

func (s *LicenseTestSuite) Test_GetLicenseOfv010() {
	checker, _ := NewChecker(context.TODO())

	latestVersion, err := checker.Type("github.com/beatlabs/patron", "v0.1.0")

	s.NoError(err)
	s.Equal("Apache-2.0", latestVersion)
}

func (s *LicenseTestSuite) Test_EscapeForUpperCase() {
	encodedName := encodeModuleName("github.com/Azure/azure-sdk-for-go")

	s.Equal("github.com/!azure/azure-sdk-for-go", encodedName)
}

func (s *LicenseTestSuite) Test_EscapeForLowerCase() {
	encodedName := encodeModuleName("github.com/beatlabs/patron")

	s.Equal("github.com/beatlabs/patron", encodedName)
}

func (s *LicenseTestSuite) Test_AllDependenciesLicense() {
	checker, _ := NewChecker(context.TODO())

	types, err := checker.Types(s.tempDir)

	s.NoError(err)
	s.NotEmpty(types)

	for name, result := range types {
		s.NoError(result.Error, name)
	}
}
