package print

import (
	"fmt"
	"heimatcli/src/heimat"
	"time"

	"github.com/alexeyco/simpletable"
)

func DescStats(stats heimat.DescStats, start, end time.Time) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Project"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Desc"},
			{Align: simpletable.AlignCenter, Text: "Time Spent [hh:mm:ss]"},
			{Align: simpletable.AlignCenter, Text: "Time Spent [%]"},
		},
	}
	sorted := stats.Sorted()
	for _, timeSpent := range sorted {

		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: timeSpent.DescTimeSpent.ProjectName},
			{Align: simpletable.AlignLeft, Text: timeSpent.DescTimeSpent.TaskName},
			{Align: simpletable.AlignLeft, Text: timeSpent.Desc},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", timeSpent.DescTimeSpent.Absolute)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%f", timeSpent.DescTimeSpent.Relative)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	fmt.Printf(" %s - %s \n\n", FormatDate(start), FormatDate(end))
	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())
}
