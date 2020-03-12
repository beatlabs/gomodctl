# gomodctl

[![gomodctl build](https://github.com/beatlabs/gomodctl/workflows/gomodctl%20build/badge.svg)](https://github.com/beatlabs/gomodctl/actions?query=workflow%3A%22gomodctl+build%22)
[![gomodctl release](https://github.com/beatlabs/gomodctl/workflows/gomodctl%20release/badge.svg)](https://github.com/beatlabs/gomodctl/actions?query=workflow%3A%22gomodctl+release%22)
[![gomodctl snyk](https://github.com/beatlabs/gomodctl/workflows/gomodctl%20snyk/badge.svg)](https://github.com/beatlabs/gomodctl/actions?query=workflow%3A%22gomodctl+snyk%22)

*gomodctl* - search, check and update go modules.

Currently supported commands:
- search -  search for Go packages by the given term
- info - search by the given term and show information about the matched package
- check - check project dependencies for the version information and shows outdated packages
- update - automatically sync project dependencies with their latest version
- license - fetch license of a module with/without version

## Installation

Execute:

```bash
$ go get github.com/beatlabs/gomodctl
```

Or using [Homebrew 🍺](https://brew.sh)

```bash
brew tap beatlabs/gomodctl https://github.com/beatlabs/gomodctl
brew install gomodctl
```


## Features

### gomodctl search <term>

Search in Go registry and return matched results.

Command:

```shell script
gomodctl search patron
```

Result:

```shell script
                   NAME                   | STARS |   IMPORT COUNT    |  SCORE
------------------------------------------+-------+-------------------+-----------
  github.com/beatlabs/patron/log          |    44 |                26 | 0.940500
  github.com/beatlabs/patron/trace        |    43 |                19 | 0.940500
  github.com/mantzas/patron/log           |    26 |                12 | 0.940500
  github.com/mantzas/patron/errors        |     9 |                12 | 0.940500
  github.com/beatlabs/patron              |    44 |                 9 | 1.000000
  github.com/beatlabs/patron/async        |    43 |                10 | 0.940500
  github.com/mantzas/patron/trace         |    26 |                 8 | 0.940500
  github.com/beatlabs/patron/examples     |    44 |                 8 | 0.940500
  github.com/beatlabs/patron/sync         |    46 |                 6 | 0.940500
  github.com/beatlabs/patron/sync/http    |    46 |                 6 | 0.931095
  github.com/mantzas/patron/encoding      |    20 |                 5 | 0.940500
  github.com/mantzas/patron               |    24 |                 4 | 1.000000
  github.com/mantzas/patron/encoding/json |     9 |                 5 | 0.931095
  github.com/mantzas/patron/examples      |    26 |                 4 | 0.940500
  github.com/mantzas/patron/async         |    24 |                 4 | 0.940500
  github.com/beatlabs/patron/trace/http   |    46 |                 4 | 0.931095
  github.com/mantzas/patron/sync          |    26 |                 3 | 0.940500
  github.com/mantzas/patron/sync/http     |    26 |                 3 | 0.931095
  github.com/mantzas/patron/info          |    13 |                 2 | 0.940500
  github.com/beatlabs/patron/async/kafka  |    44 |                 2 | 0.931095
------------------------------------------+-------+-------------------+-----------
                                                    NUMBER OF MODULES |    20
                                                  --------------------+-----------
```

Add `--json` parameter to the command to print result as a JSON.

Command:

```shell script
gomodctl search github --json
```

Result:

```json
...
{
  "Name": "gomock",
  "Path": "github.com/golang/mock/gomock",
  "ImportCount": 3900,
  "Stars": 3844,
  "Score": 0.99,
  "Synopsis": "Package gomock is a mock framework for Go."
},
{
  "Name": "uuid",
  "Path": "github.com/pborman/uuid",
  "ImportCount": 3579,
  "Stars": 318,
  "Score": 1,
  "Synopsis": "The uuid package generates and inspects UUIDs."
},
{
  "Name": "distribution",
  "Path": "github.com/docker/distribution",
  "ImportCount": 3549,
  "Stars": 5219,
  "Score": 1,
  "Synopsis": "Package distribution will define the interfaces for the components of docker distribution."
},
{
  "Name": "state",
  "Path": "github.com/juju/juju/state",
  "ImportCount": 3851,
  "Stars": 1614,
  "Score": 0.99,
  "Synopsis": ""
},
...
```

### gomodctl info <term>

Detailed information about the package with fetched documentation.

Command:

```shell script
gomodctl info patron
```

Result:

```shell script
                PATH               | STARS | IMPORT COUNT |  SCORE
-----------------------------------+-------+--------------+-----------
  github.com/beatlabs/patron/trace |    43 |           19 | 0.940500

Documentation:
PACKAGE

package trace
    import "github.com/beatlabs/patron/trace"
...
```

### gomodctl check

Check module versions in the given Go project.

Command:

```shell script
gomodctl check
```

Result:

```shell script
              MODULE              |       CURRENT       |       LATEST
----------------------------------+---------------------+----------------------
  github.com/stretchr/testify     | v1.3.0              | v1.4.0
  go.mongodb.org/mongo-driver     | v1.1.1              | v1.2.1
  github.com/mitchellh/go-homedir | v1.1.0              | v1.1.0
  github.com/ory/dockertest       | v3.3.5+incompatible | v3.3.5+incompatible
  github.com/pkg/errors           | v0.8.1              | v0.9.1
  github.com/spf13/cobra          | v0.0.5              | v0.0.5
  github.com/spf13/viper          | v1.4.0              | v1.6.2
----------------------------------+---------------------+----------------------
                                     NUMBER OF MODULES  |          7
                                  ----------------------+----------------------
```

Add `--json` parameter to the command to print result as a JSON.

### gomodctl scan

Scan for vulnerabilities using the tool [gosec](https://github.com/securego/gosec)

Command:

```shell script
gomodctl scan
```

Result
```shell script
                MODULE                | CONFIDENCE | SEVERITY |                       CWE                       |                                                                    LINE,COLUMN
--------------------------------------+------------+----------+-------------------------------------------------+------------------------------------------------------------------------------------
  github.com/mitchellh/go-homedir     | HIGH       | MEDIUM   | https://cwe.mitre.org/data/definitions/78.html  | /go/pkg/mod/github.com/mitchellh/go-homedir@v1.1.0/homedir.go
                                      |            |          |                                                 | ln:114 | col:10  exec.Command("getent", "passwd", strconv.Itoa(os.Getuid()))
--------------------------------------+------------+----------+-------------------------------------------------+------------------------------------------------------------------------------------
  github.com/olekukonko/tablewriter   | HIGH       | MEDIUM   | https://cwe.mitre.org/data/definitions/22.html  | /go/pkg/mod/github.com/olekukonko/tablewriter@v0.0.4/csv.go
                                      |            |          |                                                 | ln:19 | col:15  os.Open(fileName)
--------------------------------------+------------+----------+-------------------------------------------------+------------------------------------------------------------------------------------
  github.com/olekukonko/tablewriter   | HIGH       | LOW      | https://cwe.mitre.org/data/definitions/703.html | /go/pkg/mod/github.com/olekukonko/tablewriter@v0.0.4/table.go
                                      |            |          |                                                 | ln:782 | col:3  tmpWriter.WriteTo(t.out)
--------------------------------------+------------+----------+-------------------------------------------------+------------------------------------------------------------------------------------
  github.com/stretchr/testify         | HIGH       | MEDIUM   | https://cwe.mitre.org/data/definitions/22.html  | /Users/thomaschavakis/go/pkg/mod/github.com/stretchr/testify@v1.4.0/_codegen/main.go
                                      |            |          |                                                 | ln:184 | col:15  ioutil.ReadFile(path.Join(pd.SrcRoot, pd.ImportPath, fname))
--------------------------------------+------------+----------+-------------------------------------------------+------------------------------------------------------------------------------------
```

### gomodctl update

Update module versions to latest minor

Command:

```shell script
gomodctl update
```

Result:

```shell script
Your dependencies updated to latest minor and go.mod.backup created
              MODULE              |      PREVIOUS       |         NOW
----------------------------------+---------------------+----------------------
  github.com/ory/dockertest       | v3.3.5+incompatible | v3.3.5+incompatible
  github.com/pkg/errors           | v0.8.1              | v0.9.1
  github.com/spf13/cobra          | v0.0.5              | v0.0.5
  github.com/spf13/viper          | v1.4.0              | v1.6.2
  github.com/stretchr/testify     | v1.3.0              | v1.4.0
  go.mongodb.org/mongo-driver     | v1.1.1              | v1.2.1
  github.com/mitchellh/go-homedir | v1.1.0              | v1.1.0
----------------------------------+---------------------+----------------------
                                     NUMBER OF MODULES  |          7
                                  ----------------------+----------------------
```

### gomodctl license <modulename> <version>

Fetch license of a module, version is optional

Command:

```shell script
gomodctl license github.com/beatlabs/patron
```

Result:

```shell script
Apache-2.0
```

## How to configure for private modules

Since check and update rely on go toolchain, if you have any private module that isn't publicly accessible, don't forget to set up your environment variables. For more information and how to configure, please check [Module configuration for non-public modules](https://golang.org/cmd/go/#hdr-Module_configuration_for_non_public_modules).

## Code of conduct

Please note that this project is released with a [Contributor Code of Conduct](https://github.com/beatlabs/gomodctl/blob/master/CODE_OF_CONDUCT.md). By participating in this project and its community you agree to abide by those terms.
