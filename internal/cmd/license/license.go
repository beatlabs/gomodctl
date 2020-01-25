package license

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// Typer is exported.
type Typer interface {
	Type(moduleName, version string) (string, error)
}

// Options is exported.
type Options struct {
	Module  string
	Version string
}

// NewCmdLicense returns an instance of License command.
func NewCmdLicense(typer Typer) *cobra.Command {
	o := Options{}

	cmd := &cobra.Command{
		Use:   "license [module name] [version]",
		Short: "fetch license of module, version is optional",
		Long:  `fetch license of module, if version is empty it will use latest version`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires module name to search")
			}

			o.Module = args[0]
			if len(args) > 1 {
				o.Version = args[1]
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Execute(typer)
		},
	}

	return cmd
}

// Execute is exported.
func (o *Options) Execute(op Typer) {
	licenseType, err := op.Type(o.Module, o.Version)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(licenseType)
}
