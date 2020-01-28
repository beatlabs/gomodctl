package search

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

const (
	// LIMIT output results by default.
	LIMIT = 20
)

// Searcher is exported.
type Searcher interface {
	Search(term string) ([]internal.SearchResult, error)
}

// Options is exported.
type Options struct {
	Term    string
	ShowAll bool
}

// NewCmdSearch returns an instance of Search command.
func NewCmdSearch(searcher Searcher) *cobra.Command {
	o := Options{}

	cmd := &cobra.Command{
		Use:   "search [term to search]",
		Short: "search in packages",
		Long:  `search existing packages by the given term`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a term to search")
			}

			o.Term = strings.Join(args, " ")
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Fill(cmd)
			o.Execute(searcher)
		},
	}

	cmd.Flags().BoolP("show-all", "a", o.ShowAll, "--show-all")

	return cmd
}

// Fill fills flags into options.
func (o *Options) Fill(cmd *cobra.Command) {
	o.ShowAll, _ = cmd.Flags().GetBool("show-all")
}

// Execute is exported.
func (o *Options) Execute(op Searcher) {
	searchResults, err := op.Search(o.Term)
	if err != nil {
		fmt.Println(err)
	}

	if len(searchResults) == 0 {
		fmt.Printf("No match found for search term \"%s\"\n", o.Term)
		return
	}

	var data [][]string
	limit := o.calcLimit(searchResults)
	limitedResults := searchResults[:limit]

	for _, result := range limitedResults {
		data = append(data, []string{
			result.Path,
			strconv.Itoa(result.Stars),
			strconv.Itoa(result.ImportCount),
			fmt.Sprintf("%f", result.Score),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Stars", "Import count", "Score"})
	table.SetFooter([]string{"", "", "number of modules", strconv.Itoa(len(limitedResults))})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	fmt.Println()
}

func (o *Options) calcLimit(srs []internal.SearchResult) int {
	c := len(srs)
	if c > LIMIT {
		if o.ShowAll {
			return c
		}
		return LIMIT
	}
	return c
}
