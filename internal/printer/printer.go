package printer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type TableData struct {
	Header       []string
	Footer       []string
	RowSeparator string
	ShowBorder   bool
	ShowRowLine  bool
	Data         [][]string
}

type Printable interface {
	TableData() *TableData
	JsonData() interface{}
}

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

func PrintJSON(p Printable) {
	dataB, _ := json.Marshal(p)
	if dataB == nil {
		fmt.Println("no data")
		return
	}
	fmt.Println(string(dataB))
}
