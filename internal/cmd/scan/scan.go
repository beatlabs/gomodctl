package scan

import (
	"fmt"
	"os"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

// Scanner is exported.
type Scanner interface {
	Scan(path string) (map[string]internal.VulnerabilityResult, error)
}

// Options is exported.
type Options struct {
	Path string
}

// NewCmdScan returns an instance of Scan command.
func NewCmdScan(scanner Scanner) *cobra.Command {
	o := Options{}

	cmd := &cobra.Command{
		Use:   "scan [module name] [OPTIONS]",
		Short: "scan local for security vulnerabilities",
		Long:  `scan local module for security vulnerabilities`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				o.Path = ""
			} else {
				o.Path = args[0]
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Execute(scanner)
		},
	}

	return cmd
}

// Execute is exported.
func (o *Options) Execute(scanner Scanner) {
	var err error
	var vulnerabilitiesResult map[string]internal.VulnerabilityResult
	vulnerabilitiesResult, err = scanner.Scan(o.Path)
	if err != nil {
		fmt.Println(err)
		return
	}
	RenderResults(vulnerabilitiesResult)
}

// RenderResults renders the vulnerabilities.
func RenderResults(vulnerabilitiesResult map[string]internal.VulnerabilityResult) {
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
