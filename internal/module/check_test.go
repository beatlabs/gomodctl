package module

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CheckTestSuite struct {
	suite.Suite
}

func TestCheckTestSuite(t *testing.T) {
	suite.Run(t, new(CheckTestSuite))
}

func (s *CheckTestSuite) Test_SelfCheck() {
	checker := Checker{}

	result, err := checker.Check("../..")
	s.NoError(err)
	s.Len(result, 7)
}
