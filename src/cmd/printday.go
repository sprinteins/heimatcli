package main

import (
	"fmt"
	"heimatcli/src/heimat"
	"time"

	"github.com/alexeyco/simpletable"
)

func printDay(day *heimat.Day) {

	printHeimatDate(day.Date)

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
	for _, tt := range day.TrackedTimes {
		dur := calcDuration(tt.Start, tt.End)
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
