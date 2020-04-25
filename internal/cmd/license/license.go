package license

import (
	"fmt"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/beatlabs/gomodctl/internal/printer"
	"github.com/spf13/cobra"
)

// Typer defines interface to check for license types.
type Typer interface {
	Type(moduleName, version string) (string, error)
	Types(path string) (map[string]internal.LicenseResult, error)
}

// Options contains module and version to check.
// Both are optional.
type Options struct {
	Module  string
	Version string
	JSON    bool
	Path    string
}

// NewCmdLicense returns an instance of License command.
func NewCmdLicense(typer Typer) *cobra.Command {
	o := Options{}

	cmd := &cobra.Command{
		Use:   "license [module name] [version]",
		Short: "fetch license of module, version is optional",
		Long:  `fetch license of module, if version is empty it will use latest version`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				o.Module = args[0]
			}

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

// Fill fills flags into options.
func (o *Options) Fill(cmd *cobra.Command) {
	o.JSON, _ = cmd.Flags().GetBool("json")
	o.Path, _ = cmd.Flags().GetString("path")
}

// Execute executes command on given Typer and prints output.
func (o *Options) Execute(op Typer) {
	if o.Version == "" && o.Module == "" {
		types, err := op.Types(o.Path)
		if err != nil {
			fmt.Println(err)
			return
		}

		rp := NewResultPrinter(types)
		if o.JSON {
			printer.PrintJSON(rp)
		} else {
			printer.PrintTable(rp)
		}
	} else {
		licenseType, err := op.Type(o.Module, o.Version)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(licenseType)
	}
}
