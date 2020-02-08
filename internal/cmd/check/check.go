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
	VulnerabilitiesCheck(path string, vulnerabilityCheck bool, jsonOutputCheck bool) (map[string]internal.VulnerabilityResult, error)
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
		Use:   "check [module name] [OPTIONS]",
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
			f, _ := cmd.Flags().GetBool("vulnerabilities")
			j, _ := cmd.Flags().GetBool("json")
			o.Execute(checker, f, j)
		},
	}

	cmd.Flags().BoolP("vulnerabilities", "v", false, "Check for vulnerabilities")
	cmd.Flags().BoolP("json", "j", false, "output to a json file")

	return cmd
}

// Execute is exported.
func (o *Options) Execute(checker Checker, vulnerabilitiesCheck bool, jsonOutputCheck bool) {
	var checkResults map[string]internal.CheckResult
	var err error
	var vulnerabilitiesResult map[string]internal.VulnerabilityResult

	if vulnerabilitiesCheck {
		vulnerabilitiesResult, err = checker.VulnerabilitiesCheck(o.Path, vulnerabilitiesCheck, jsonOutputCheck)
		if err != nil {
			fmt.Println(err)
			return
		}
		vulnerabilitiesResultsRender(vulnerabilitiesResult)
	} else {
		checkResults, err = checker.Check(o.Path)
		if err != nil {
			fmt.Println(err)
			return
		}
		checkResultsRender(checkResults)
	}
}

func vulnerabilitiesResultsRender(vulnerabilitiesResult map[string]internal.VulnerabilityResult) {
	var data [][]string
	for name, result := range vulnerabilitiesResult {
		for _, issue := range result.Issues {
			data = append(data, []string{
				name,
				issue.Confidence,
				issue.Severity,
				issue.Cwe.URL,
				fmt.Sprintf("%s\nln:%s | col:%s \n%s", issue.File, issue.Line, issue.Column, issue.Code),
			})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module", "Confidence", "Severity", "CWE", "Line,Column"})
	table.SetBorder(false)
	table.SetRowLine(true)
	table.SetRowSeparator("-")
	table.AppendBulk(data)
	table.Render()
}

func checkResultsRender(checkResults map[string]internal.CheckResult) {

	var data [][]string

	for name, result := range checkResults {

		var localVersion = ""
		if result.LocalVersion != nil {
			localVersion = result.LocalVersion.Original()
		}
		var latestVersion = ""
		if result.LatestVersion != nil {
			latestVersion = result.LatestVersion.Original()
		}

		data = append(data, []string{
			name,
			localVersion,
			latestVersion,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module", "Current", "Latest"})
	table.SetFooter([]string{"", "number of modules", strconv.Itoa(len(checkResults))})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
