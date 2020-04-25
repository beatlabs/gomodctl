package license

import (
	"fmt"
	"strconv"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/beatlabs/gomodctl/internal/printer"
)

// ResultPrinter implements Printer interface for License command.
type ResultPrinter struct {
	licenseResults map[string]internal.LicenseResult
}

// NewResultPrinter creates a new instance of ResultPrinter.
func NewResultPrinter(m map[string]internal.LicenseResult) *ResultPrinter {
	return &ResultPrinter{licenseResults: m}
}

// TableData returns table friendly result.
func (r *ResultPrinter) TableData() *printer.TableData {
	var data [][]string

	for name, result := range r.licenseResults {
		r := []string{
			name,
			result.LocalVersion.Original(),
		}

		if result.Error != nil {
			r = append(r, fmt.Sprintf("failed because of: %s", result.Error.Error()))
		} else {
			r = append(r, result.Type)
		}

		data = append(data, r)
	}

	td := &printer.TableData{
		Header:       []string{"Module", "Version", "License"},
		Footer:       []string{"", "number of modules", strconv.Itoa(len(r.licenseResults))},
		RowSeparator: "-",
		ShowBorder:   false,
		ShowRowLine:  false,
		Data:         data,
	}

	return td
}

// JSONData returns JSON friendly result.
func (r *ResultPrinter) JSONData() interface{} {
	return r.licenseResults
}
