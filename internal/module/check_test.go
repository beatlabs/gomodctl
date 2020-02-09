package module

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
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
func (s *CheckTestSuite) Test_CustomModFile() {
	gomodContent := []byte(`module test-project

go 1.13

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-gonic/gin v1.5.0 // indirect
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-openapi/spec v0.19.4 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/klauspost/compress v1.9.3 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/labstack/echo/v4 v4.1.11
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/newrelic/go-agent v2.16.0+incompatible // indirect
	github.com/prometheus/client_golang v1.2.1 // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/swaggo/echo-swagger v0.0.0-20190329130007-1219b460a043
	github.com/swaggo/swag v1.6.3 // indirect
	github.com/valyala/fasthttp v1.6.0 // indirect
	github.com/valyala/fasttemplate v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20191128160524-b544559bb6d1 // indirect
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933
	golang.org/x/sys v0.0.0-20191128015809-6d18c012aee9 // indirect
	golang.org/x/tools v0.0.0-20191130070609-6e064ea0cf2d // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/validator.v9 v9.30.2 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)`)

	tempDir, err := ioutil.TempDir("", "test")
	s.NoError(err)

	tempFile := filepath.Join(tempDir, "go.mod")

	err = ioutil.WriteFile(tempFile, gomodContent, 0666)
	s.NoError(err)

	checker := Checker{Ctx: s.ctx}

	result, err := checker.Check(tempDir)
	s.NoError(err)
	s.NotEmpty(result)

	s.NoError(os.Remove(tempFile))
	s.NoError(os.RemoveAll(tempDir))
}
