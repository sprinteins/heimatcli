package print

import (
	"fmt"
	"heimatcli/src/heimat"
	"time"

	"github.com/alexeyco/simpletable"
)

func ProjectStats(stats heimat.ProjectStats, start, end time.Time) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Project"},
			{Align: simpletable.AlignCenter, Text: "Time Spent [hh:mm:ss]"},
			{Align: simpletable.AlignCenter, Text: "Time Spent [%]"},
		},
	}

	sorted := stats.Sorted()
	for _, timeSpent := range sorted {

		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: timeSpent.ProjectName},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", timeSpent.ProjectTimeSpent.Absolute)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%f", timeSpent.ProjectTimeSpent.Relative)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Printf(" %s - %s \n\n", FormatDate(start), FormatDate(end))
	fmt.Println(table.String())
}
