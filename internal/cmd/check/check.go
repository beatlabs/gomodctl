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
	VulnerabilitiesCheck(path string, vulnerabilityCheck bool) (map[string]internal.VulnerabilityResult, error)
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
		Use:   "check",
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
			o.Execute(checker, f)
		},
	}

	cmd.Flags().BoolP("vulnerabilities", "v", false, "Check for vulnerabilities")

	return cmd
}

// Execute is exported.
func (o *Options) Execute(checker Checker, vulnerabilitiesCheck bool) {
	var checkResults map[string]internal.CheckResult
	var err error
	var vulnerabilitiesResult map[string]internal.VulnerabilityResult

	if vulnerabilitiesCheck {
		vulnerabilitiesResult, err = checker.VulnerabilitiesCheck(o.Path, vulnerabilitiesCheck)
		if err != nil {
			fmt.Println(err)
			return
		}
		vulnerabilitesResultsRender(vulnerabilitiesResult)
	} else {
		checkResults, err = checker.Check(o.Path)
		if err != nil {
			fmt.Println(err)
			return
		}
		checkResultsRender(checkResults)
	}
}

func vulnerabilitesResultsRender(vulnerabilitiesResult map[string]internal.VulnerabilityResult) {
	var data [][]string
	for name, result := range vulnerabilitiesResult {
		fmt.Println(name)
		for _, issue := range result.Issues {
			data = append(data, []string{
				name,
				issue.Confidence,
				issue.Severity,
				issue.Cwe.URL,
				issue.Code,
				fmt.Sprintf("%s,%s", issue.Line, issue.Column),
				issue.File,
			})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module", "Confidence", "Severity", "CWE", "Code", "Line,Column", "File"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}

func checkResultsRender(checkResults map[string]internal.CheckResult) {

	var data [][]string

	for name, result := range checkResults {
		data = append(data, []string{
			name,
			result.LocalVersion.Original(),
			result.LatestVersion.Original(),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module", "Current", "Latest"})
	table.SetFooter([]string{"", "number of modules", strconv.Itoa(len(checkResults))})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
