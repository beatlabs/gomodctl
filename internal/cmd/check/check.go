package check

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/beatlabs/gomodctl/internal"
)

// Checker is exported.
type Checker interface {
	Check(path string) ([]internal.CheckResult, error)
}

// CheckOptions is exported.
type CheckOptions struct {
	Path string
}

// NewCmdCheck returns an instance of Search command.
func NewCmdCheck(checker Checker) *cobra.Command {
	o := CheckOptions{}

	cmd := &cobra.Command{
		Use:   "check",
		Short: "check local packages for updates",
		Long:  `get list of local packages and check them for updates`,
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
func (o *CheckOptions) Execute(checker Checker) {
	checkResults, err := checker.Check(o.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data [][]string

	for _, result := range checkResults {
		data = append(data, []string{
			result.Name,
			result.LocalVersion,
			result.LatestVersion,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Package", "Current", "Latest"})
	table.SetFooter([]string{"", "number of packages", strconv.Itoa(len(checkResults))})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
