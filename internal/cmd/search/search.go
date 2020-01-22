package search

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/beatlabs/gomodctl/internal"
)

// Searcher is exported.
type Searcher interface {
	Search(term string) ([]internal.SearchResult, error)
}

// SearchOptions is exported.
type SearchOptions struct {
	Term string
}

// NewCmdSearch returns an instance of Search command.
func NewCmdSearch(searcher Searcher) *cobra.Command {
	o := SearchOptions{}

	cmd := &cobra.Command{
		Use:   "search [term to search]",
		Short: "search in packages",
		Long:  `search existing packages by the given term`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a term to search")
			}

			o.Term = args[0]
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// o.Fill(cmd)
			o.Execute(searcher)
		},
	}

	return cmd
}

// Execute is exported.
func (o *SearchOptions) Execute(op Searcher) {
	searchResults, err := op.Search(o.Term)
	if err != nil {
		fmt.Println(err)
	}

	if len(searchResults) == 0 {
		fmt.Println("No match found")
		return
	}

	var data [][]string

	for _, result := range searchResults {
		data = append(data, []string{
			result.Path,
			strconv.Itoa(result.Stars),
			strconv.Itoa(result.ImportCount),
			fmt.Sprintf("%f", result.Score),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Stars", "Import count", "Score"})
	table.SetFooter([]string{"", "", "number of packages", strconv.Itoa(len(searchResults))})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
