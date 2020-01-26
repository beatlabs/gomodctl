package module

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CheckTestSuite struct {
	suite.Suite
	cnl context.CancelFunc
	ctx context.Context
}

func TestCheckTestSuite(t *testing.T) {
	suite.Run(t, new(CheckTestSuite))
}

func (s *CheckTestSuite) SetupTest() {
	s.ctx, s.cnl = context.WithCancel(context.Background())
}

func (s *CheckTestSuite) TearDownTest() {
	s.cnl()
}

func (s *CheckTestSuite) Test_SelfCheck() {
	checker := Checker{Ctx: s.ctx}

	result, err := checker.Check("../..")
	s.NoError(err)
	s.NotEmpty(result)
}

func (s *CheckTestSuite) Test_SelfCheck_CancelContextBefore() {
	checker := Checker{Ctx: s.ctx}
	s.cnl()

	result, err := checker.Check("../..")
	s.EqualError(err, "context canceled")
	s.Empty(result)
}
