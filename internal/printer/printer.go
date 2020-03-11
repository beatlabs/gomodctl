package printer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// TableData defines print options for table output.
type TableData struct {
	Header       []string
	Footer       []string
	RowSeparator string
	ShowBorder   bool
	ShowRowLine  bool
	Data         [][]string
}

// Printable defines contract to be implemented in order to print
// commands output as a JSON or Table.
type Printable interface {
	TableData() *TableData
	JSONData() interface{}
}

// PrintTable prints printable result as a table output.
func PrintTable(p Printable) {
	td := p.TableData()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(td.Header)
	table.SetFooter(td.Footer)
	table.SetRowSeparator(td.RowSeparator)
	table.SetBorder(td.ShowBorder)
	table.SetRowLine(td.ShowRowLine)
	table.AppendBulk(td.Data)
	table.Render()
}

// PrintJSON prints printable result as a JSON output.
func PrintJSON(p Printable) {
	data := p.JSONData()
	if data == nil {
		fmt.Println("no data")
		return
	}

	dataB, err := json.Marshal(data)
	if err != nil {
		fmt.Println("failed to parse json", err)
	} else {
		fmt.Println(string(dataB))
	}
}
