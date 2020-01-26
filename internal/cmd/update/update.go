package check

import (
	"fmt"
	"os"
	"strconv"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Updater is exported.
type Updater interface {
	Update(path string, updateType internal.UpdateType) (map[string]internal.CheckResult, error)
}

// Options is exported.
type Options struct {
	Path       string
	UpdateType internal.UpdateType
}

// NewCmdUpdate returns an instance of Update command.
func NewCmdUpdate(updater Updater) *cobra.Command {
	o := Options{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update project dependencies",
		Long:  `update project dependencies to minor or major versions`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				o.Path = ""
			} else {
				o.Path = args[0]
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Execute(updater)
		},
	}

	return cmd
}

// Execute is exported.
func (o *Options) Execute(updater Updater) {
	checkResults, err := updater.Update(o.Path, o.UpdateType)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data [][]string

	for name, result := range checkResults {
		data = append(data, []string{
			name,
			result.LocalVersion.Original(),
			result.LatestVersion.Original(),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module", "Previous", "Now"})
	table.SetFooter([]string{"", "number of packages", strconv.Itoa(len(checkResults))})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
