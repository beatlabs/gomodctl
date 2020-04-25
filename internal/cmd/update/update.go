package check

import (
	"fmt"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/beatlabs/gomodctl/internal/printer"
	"github.com/spf13/cobra"
)

// Updater is exported.
type Updater interface {
	Update(path string) (map[string]internal.CheckResult, error)
}

// Options is exported.
type Options struct {
	Path string
	JSON bool
}

// NewCmdUpdate returns an instance of Update command.
func NewCmdUpdate(updater Updater) *cobra.Command {
	o := Options{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update project dependencies",
		Long:  `update project dependencies to minor versions`,
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Execute(updater)
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
func (o *Options) Execute(updater Updater) {
	checkResults, err := updater.Update(o.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Your dependencies updated to latest minor and go.mod.backup created")

	rp := NewResultPrinter(checkResults)
	if o.JSON {
		printer.PrintJSON(rp)
	} else {
		printer.PrintTable(rp)
	}
}
