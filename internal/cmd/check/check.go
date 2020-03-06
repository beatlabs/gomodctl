package check

import (
	"fmt"
	"os"
	"strconv"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Checker is exported.
type Checker interface {
	Check(path string) (map[string]internal.CheckResult, error)
}

// Options is exported.
type Options struct {
	Path string
}

// NewCmdCheck returns an instance of Search command.
func NewCmdCheck(checker Checker) *cobra.Command {
	o := Options{}

	cmd := &cobra.Command{
		Use:   "check [module name]",
		Short: "check local module for updates",
		Long:  `get list of local module and check them for updates`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				o.Path = ""
			} else {
				o.Path = args[0]
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Execute(checker)
		},
	}

	return cmd
}

// Execute is exported.
func (o *Options) Execute(checker Checker) {
	checkResults, err := checker.Check(o.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data [][]string

	for name, result := range checkResults {
		r := []string{
			name,
			result.LocalVersion.Original(),
		}

		if result.Error != nil {
			r = append(r, result.Error.Error())
		} else {
			r = append(r, result.LatestVersion.Original())
		}

		data = append(data, r)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module", "Current", "Latest"})
	table.SetFooter([]string{"", "number of modules", strconv.Itoa(len(checkResults))})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
