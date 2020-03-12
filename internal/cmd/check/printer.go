package check

import (
	"strconv"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/beatlabs/gomodctl/internal/printer"
)

// ResultPrinter implements Printer interface for Check command.
type ResultPrinter struct {
	Result map[string]internal.CheckResult
}

// NewResultPrinter creates a new instance of ResultPrinter.
func NewResultPrinter(results map[string]internal.CheckResult) *ResultPrinter {
	return &ResultPrinter{
		Result: results,
	}
}

// TableData returns table friendly result.
func (p *ResultPrinter) TableData() *printer.TableData {
	var data [][]string

	for name, result := range p.Result {
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

	td := &printer.TableData{
		Header:       []string{"Module", "Current", "Latest"},
		Footer:       []string{"", "number of modules", strconv.Itoa(len(p.Result))},
		RowSeparator: "-",
		ShowBorder:   false,
		ShowRowLine:  false,
		Data:         data,
	}

	return td
}

// JSONData returns JSON friendly result.
func (p *ResultPrinter) JSONData() interface{} {
	return p.Result
}
