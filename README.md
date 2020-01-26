# gomodctl [![actions](https://github.com/beatlabs/gomodctl/workflows/gomodctl%20build/badge.svg)](https://github.com/beatlabs/gomodctl/actions)

*gomodctl* - search, check and update go modules.

Currently supported commands:
- search -  search go packages (godoc) by the given search term
- info - print information about the first matched package (eg `gomodctl info gomock`)
- check - check packages in the project for new versions (for now it requires go.mod file)

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
                          NAME                          | STARS |    IMPORT COUNT    |  SCORE
--------------------------------------------------------+-------+--------------------+-----------
  github.com/beatlabs/patron/trace                      |    43 |                 19 | 0.940500
  github.com/mantzas/patron/log                         |    26 |                 12 | 0.940500
  github.com/mantzas/patron/errors                      |     9 |                 12 | 0.940500
  github.com/beatlabs/patron                            |    43 |                  9 | 1.000000
  github.com/beatlabs/patron/async                      |    43 |                 10 | 0.940500
  github.com/beatlabs/patron/examples                   |    38 |                  9 | 0.940500
  github.com/mantzas/patron/trace                       |    26 |                  8 | 0.940500
  github.com/beatlabs/patron/sync                       |    43 |                  6 | 0.940500
  github.com/beatlabs/patron/sync/http                  |    41 |                  6 | 0.931095
  github.com/mantzas/patron/encoding                    |    20 |                  5 | 0.940500
  github.com/mantzas/patron                             |    24 |                  4 | 1.000000
  github.com/mantzas/patron/encoding/json               |     9 |                  5 | 0.931095
  github.com/mantzas/patron/examples                    |    26 |                  4 | 0.940500
  github.com/mantzas/patron/async                       |    24 |                  4 | 0.940500
  github.com/mantzas/patron/sync                        |    26 |                  3 | 0.940500
  github.com/mantzas/patron/sync/http                   |    26 |                  3 | 0.931095
  github.com/mantzas/patron/info                        |    13 |                  2 | 0.940500
  github.com/beatlabs/patron/async/kafka                |    33 |                  2 | 0.931095
  github.com/beatlabs/patron/trace/amqp                 |    43 |                  2 | 0.931095
  github.com/mantzas/patron/trace/http                  |    24 |                  2 | 0.931095
  github.com/beatlabs/patron/async/amqp                 |    36 |                  2 | 0.931095
  github.com/beatlabs/patron/trace/kafka                |    36 |                  2 | 0.931095
  github.com/beatlabs/patron/log                        |    19 |                  1 | 0.940500
  github.com/beatlabs/patron/correlation                |    43 |                  1 | 0.940500
  github.com/mantzas/patron/async/kafka                 |    26 |                  1 | 0.931095
  github.com/mantzas/patron/async/amqp                  |    24 |                  1 | 0.931095
  github.com/beatlabs/patron/log/zerolog                |    19 |                  1 | 0.931095
  github.com/mantzas/patron/log/zerolog                 |    26 |                  1 | 0.931095
  github.com/mantzas/patron/trace/amqp                  |    26 |                  1 | 0.931095
  github.com/mantzas/patron/trace/kafka                 |    26 |                  1 | 0.931095
  github.com/mantzas/patron/reliability/circuitbreaker  |    22 |                  1 | 0.931095
  github.com/beatlabs/patron/reliability/circuitbreaker |    19 |                  1 | 0.931095
  github.com/mantzas/patron/sync/http/auth              |    20 |                  1 | 0.921784
  github.com/beatlabs/patron/trace/sns                  |    36 |                  0 | 0.980100
  github.com/beatlabs/patron/trace/http                 |    19 |                  0 | 0.931095
  github.com/beatlabs/patron/trace/sql                  |    43 |                  0 | 0.931095
  github.com/beatlabs/patron/async/sqs                  |    36 |                  0 | 0.931095
  github.com/beatlabs/patron/trace/es                   |    43 |                  0 | 0.931095
  github.com/mantzas/patron/trace/sql                   |    24 |                  0 | 0.931095
  github.com/mantzas/patron/encoding/protobuf           |    20 |                  0 | 0.931095
  github.com/mantzas/patron/reliability/retry           |    24 |                  0 | 0.931095
  github.com/mantzas/patron/sync/http/auth/apikey       |    20 |                  0 | 0.912566
--------------------------------------------------------+-------+--------------------+-----------
                                                                  NUMBER OF PACKAGES |    42
                                                                ---------------------+-----------
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
```

### gomodctl check

Check package versions in the given Go project.

Command:

```shell script
gomodctl check
```

Result:

```shell script
               PACKAGE              |               CURRENT                | LATEST
------------------------------------+--------------------------------------+----------
  github.com/go-resty/resty/v2      | v2.1.0                               | v2.1.0
  github.com/mitchellh/go-homedir   | v1.1.0                               | v1.1.0
  github.com/olekukonko/tablewriter | v0.0.4                               | v0.0.4
  github.com/prometheus/common      | v0.4.0                               | 0.7.0
  github.com/spf13/cobra            | v0.0.5                               | v0.0.5
  github.com/spf13/viper            | v1.5.0                               | 1.6.0
  github.com/stretchr/testify       | v1.4.0                               | v1.4.0
  github.com/tcnksm/go-latest       | v0.0.0-20170313132115-e3007ae9052e   | 0.1.1
  golang.org/x/mod                  | v0.1.1-0.20191126161957-788aebd06792 | unknown
------------------------------------+--------------------------------------+----------
                                               NUMBER OF PACKAGES          |    9
                                    ---------------------------------------+----------
```

## Code of conduct

Please note that this project is released with a [Contributor Code of Conduct](https://github.com/beatlabs/gomodctl/blob/master/CODE_OF_CONDUCT.md). By participating in this project and its community you agree to abide by those terms.
