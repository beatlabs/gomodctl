package check

import (
	"fmt"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/beatlabs/gomodctl/internal/printer"
	"github.com/spf13/cobra"
)

// Checker is exported.
type Checker interface {
	Check(path string) (map[string]internal.CheckResult, error)
}

// Options is exported.
type Options struct {
	Path string
	JSON bool
}

// NewCmdCheck returns an instance of Search command.
func NewCmdCheck(checker Checker) *cobra.Command {
	o := Options{}

	cmd := &cobra.Command{
		Use:   "check [module name]",
		Short: "check local module for updates",
		Long:  `get list of local module and check them for updates`,
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Fill(cmd)
			o.Execute(checker)
		},
	}

	return cmd
}

// Fill fills flags into options.
func (o *Options) Fill(cmd *cobra.Command) {
	o.JSON, _ = cmd.Flags().GetBool("json")
	o.Path, _ = cmd.Flags().GetString("path")
}

// Execute is exported.
func (o *Options) Execute(checker Checker) {
	checkResults, err := checker.Check(o.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	rp := NewResultPrinter(checkResults)
	if o.JSON {
		printer.PrintJSON(rp)
	} else {
		printer.PrintTable(rp)
	}
}
