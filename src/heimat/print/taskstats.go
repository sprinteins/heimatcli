package print

import (
	"fmt"
	"heimatcli/src/heimat"
	"time"

	"github.com/alexeyco/simpletable"
)

func TaskStats(stats heimat.TaskStats, start, end time.Time) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Project"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Time Spent [hh:mm:ss]"},
			{Align: simpletable.AlignCenter, Text: "Time Spent [%]"},
		},
	}

	sorted := stats.Sorted()
	for _, timeSpent := range sorted {

		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: timeSpent.TaskTimeSpent.ProjectName},
			{Align: simpletable.AlignLeft, Text: timeSpent.TaskName},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", timeSpent.TaskTimeSpent.Absolute)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%f", timeSpent.TaskTimeSpent.Relative)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	fmt.Printf(" %s - %s \n\n", FormatDate(start), FormatDate(end))
	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())
}
