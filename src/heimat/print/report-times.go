package print

import (
	"fmt"

	"github.com/alexeyco/simpletable"
)

func TimeReports(trs []TimeReport) {

	// HeimatDate(day.Date)

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Time_Balance"},
			{Align: simpletable.AlignCenter, Text: "Unplanned_Vacations"},
		},
	}

	for _, tr := range trs {

		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: tr.Name},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%.0f", tr.TimeBalance)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%.0f", tr.VacationLeft)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())
}

type TimeReport struct {
	Name         string
	TimeBalance  float32
	VacationLeft float32
}
