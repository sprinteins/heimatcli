package print

import (
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/calc"
	"time"

	"github.com/alexeyco/simpletable"
)

func Day(day *heimat.Day) {

	HeimatDate(day.Date)

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Start"},
			{Align: simpletable.AlignCenter, Text: "End"},
			{Align: simpletable.AlignCenter, Text: "Duration"},
			{Align: simpletable.AlignCenter, Text: "Project"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Notes"},
		},
	}

	var subTime time.Duration
	for tti, tt := range day.TrackedTimes {

		prevTime := prevTime(tti, day.TrackedTimes)
		if prevTime != nil {
			diffBetweenPrevAndCurrent := calc.Duration(prevTime.End, tt.Start)
			if diffBetweenPrevAndCurrent > 0 {
				r := []*simpletable.Cell{
					{Align: simpletable.AlignRight, Text: prevTime.End},
					{Align: simpletable.AlignRight, Text: tt.Start},
					{Align: simpletable.AlignLeft, Text: diffBetweenPrevAndCurrent.String()},
					{Align: simpletable.AlignLeft, Text: "--- BREAK ---", Span: 3},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}
		}

		dur := calc.Duration(tt.Start, tt.End)
		subTime = subTime + dur

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: tt.Start},
			{Align: simpletable.AlignRight, Text: tt.End},
			{Align: simpletable.AlignLeft, Text: dur.String()},
			{Align: simpletable.AlignLeft, Text: tt.Project.Name},
			{Align: simpletable.AlignLeft, Text: tt.Task.Name},
			{Align: simpletable.AlignLeft, Text: tt.Note},
		}
		table.Body.Cells = append(table.Body.Cells, r)

	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Sum:"},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", subTime)},
			{},
			{},
			{},
			{},
		},
	}

	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())
}

func prevTime(timeIndex int, trackedTimes []heimat.TrackEntry) *heimat.TrackEntry {
	prevIndex := timeIndex - 1
	if prevIndex < 0 || prevIndex > len(trackedTimes)-1 {
		return nil
	}
	trackedTime := trackedTimes[prevIndex]
	return &trackedTime
}
