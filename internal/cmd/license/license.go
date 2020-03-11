package license

import (
	"fmt"
	"os"
	"strconv"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/olekukonko/tablewriter"
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

// Execute executes command on given Typer and prints output.
func (o *Options) Execute(op Typer) {
	if o.Version == "" && o.Module == "" {
		types, err := op.Types("")
		if err != nil {
			fmt.Println(err)
			return
		}

		var data [][]string

		for name, result := range types {
			r := []string{
				name,
				result.LocalVersion.Original(),
			}

			if result.Error != nil {
				r = append(r, fmt.Sprintf("failed because of: %s", result.Error.Error()))
			} else {
				r = append(r, result.Type)
			}

			data = append(data, r)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Module", "Version", "License"})
		table.SetFooter([]string{"", "number of modules", strconv.Itoa(len(types))})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	} else {
		licenseType, err := op.Type(o.Module, o.Version)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(licenseType)
	}
}
