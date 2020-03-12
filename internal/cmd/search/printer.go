package search

import (
	"fmt"
	"strconv"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/beatlabs/gomodctl/internal/printer"
)

// ResultPrinter implements Printer interface for Search command.
type ResultPrinter struct {
	List    []internal.SearchResult
	ShowAll bool
}

// NewResultPrinter creates a new instance of ResultPrinter.
func NewResultPrinter(results []internal.SearchResult, showAll bool) *ResultPrinter {
	return &ResultPrinter{
		List:    results,
		ShowAll: showAll,
	}
}

func (p *ResultPrinter) calcLimit(srs []internal.SearchResult) int {
	c := len(srs)
	if c > LIMIT {
		if p.ShowAll {
			return c
		}
		return LIMIT
	}
	return c
}

// TableData returns table friendly result.
func (p *ResultPrinter) TableData() *printer.TableData {
	var data [][]string
	limit := p.calcLimit(p.List)
	limitedResults := p.List[:limit]

	for _, result := range limitedResults {
		data = append(data, []string{
			result.Path,
			strconv.Itoa(result.Stars),
			strconv.Itoa(result.ImportCount),
			fmt.Sprintf("%f", result.Score),
		})
	}

	td := &printer.TableData{
		Header:       []string{"Name", "Stars", "Import count", "Score"},
		Footer:       []string{"", "", "number of modules", strconv.Itoa(len(data))},
		RowSeparator: "-",
		ShowBorder:   false,
		ShowRowLine:  false,
		Data:         data,
	}

	return td
}

// JSONData returns JSON friendly result.
func (p *ResultPrinter) JSONData() interface{} {
	return p.List
}
