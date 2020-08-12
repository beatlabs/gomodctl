package info

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Infoer is exported.
type Infoer interface {
	Search(term string) ([]internal.SearchResult, error)
	Info(path string) (string, error)
	Imports(path string) ([]string, error)
	Importers(path string) ([]string, error)
}

// Options is exported.
type Options struct {
	Term          string
	ShowImports   bool
	ShowImporters bool
}

// NewCmdInfo returns an instance of Search command.
func NewCmdInfo(ig Infoer) *cobra.Command {
	o := Options{}

	cmd := &cobra.Command{
		Use:   "info [name of the package]",
		Short: "package info",
		Long:  `return detailed info about first matched result`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a term to search")
			}

			o.Term = args[0]
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Fill(cmd)
			o.Execute(ig)
		},
	}

	cmd.Flags().BoolP("imports", "i", false, "--imports")
	cmd.Flags().BoolP("importers", "e", false, "--importers")

	return cmd
}

// Fill fills flags into options.
func (o *Options) Fill(cmd *cobra.Command) {
	o.ShowImports, _ = cmd.Flags().GetBool("imports")
	o.ShowImporters, _ = cmd.Flags().GetBool("importers")
}

// Execute is exported.
func (o *Options) Execute(ig Infoer) {
	searchResults, err := ig.Search(o.Term)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(searchResults) == 0 {
		fmt.Println("No match found")
		return
	}

	top := searchResults[0]

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Path", "Stars", "Import count", "Score"})
	table.SetBorder(false)
	table.Append([]string{
		top.Path,
		strconv.Itoa(top.Stars),
		strconv.Itoa(top.ImportCount),
		fmt.Sprintf("%f", top.Score),
	})
	table.Render()

	if o.ShowImports {
		imports, err := ig.Imports(top.Path)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("\nImports:")
		fmt.Println(strings.Join(imports, "\n"))
	}

	if o.ShowImporters {
		importers, err := ig.Importers(top.Path)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("\nImporters:")
		fmt.Println(strings.Join(importers, "\n"))
	}

	infoResult, err := ig.Info(top.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nDocumentation:")
	fmt.Println(infoResult)
}
