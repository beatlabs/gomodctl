package license

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type LicenseTestSuite struct {
	suite.Suite
}

func TestLicenseTestSuite(t *testing.T) {
	suite.Run(t, new(LicenseTestSuite))
}

func (s *LicenseTestSuite) Test_GetLatestVersion() {
	checker, _ := NewChecker()

	latestVersion, err := checker.getLatestVersion("github.com/beatlabs/patron")

	s.NotNil(latestVersion)
	s.GreaterOrEqual(latestVersion.Major(), int64(0))
	s.GreaterOrEqual(latestVersion.Minor(), int64(30))
	s.GreaterOrEqual(latestVersion.Patch(), int64(0))
	s.NoError(err)
}

func (s *LicenseTestSuite) Test_GetLatestVersionOfNonExistingPackage() {
	checker, _ := NewChecker()

	_, err := checker.getLatestVersion("github.com/beatlabs/patro")

	s.Error(err)
}

func (s *LicenseTestSuite) Test_GetLicenseOfLatest() {
	checker, _ := NewChecker()

	latestVersion, err := checker.Type("github.com/beatlabs/patron", "")

	s.Equal("Apache-2.0", latestVersion)
	s.NoError(err)
}

func (s *LicenseTestSuite) Test_GetLicenseOfv010() {
	checker, _ := NewChecker()

	latestVersion, err := checker.Type("github.com/beatlabs/patron", "v0.1.0")

	s.Equal("Apache-2.0", latestVersion)
	s.NoError(err)
}
