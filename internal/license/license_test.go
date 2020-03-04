package license

import (
	"context"
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
