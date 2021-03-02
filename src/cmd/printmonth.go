package main

import (
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/api"
	"sort"
	"time"

	"github.com/alexeyco/simpletable"
)

func printMonth(month *heimat.Month) {

	sort.Sort(byDate(month.TrackedTimesDate))
	emptyRow := []*simpletable.Cell{
		{},
		{},
		{},
		{},
	}

	d := api.DateFromHeimatDate(month.TrackedTimesDate[0].Date)
	yearMonth := d.Format("2006 01 (Jan)")
	fmt.Printf("%s\n", yearMonth)
	fmt.Printf("\n")

	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Day"},
			{Align: simpletable.AlignCenter, Text: "Time"},
			{Align: simpletable.AlignCenter, Text: "Duration"},
			{Align: simpletable.AlignCenter, Text: "Task"},
		},
	}

	for _, day := range month.TrackedTimesDate {

		d = api.DateFromHeimatDate(day.Date)
		dayDate := d.Format("02 (Mon)")

		if len(day.TrackedTimes) == 0 {
			row := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: dayDate},
				{Align: simpletable.AlignLeft, Text: "-"},
				{Align: simpletable.AlignLeft, Text: ""},
				{Align: simpletable.AlignLeft, Text: ""},
			}
			table.Body.Cells = append(table.Body.Cells, row)
			table.Body.Cells = append(table.Body.Cells, emptyRow)
			continue
		}

		var dailySum time.Duration

		for tti, trackedTime := range day.TrackedTimes {

			dur := calcDuration(trackedTime.Start, trackedTime.End)
			dailySum = dailySum + dur

			row := make([]*simpletable.Cell, 4)
			if tti == 0 {
				row[0] = &simpletable.Cell{
					Align: simpletable.AlignLeft,
					Text:  dayDate,
				}
			} else {
				row[0] = &simpletable.Cell{}
			}

			timeCell := renderTime(trackedTime)
			row[1] = &simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  timeCell,
			}

			durationCell := renderDuration(trackedTime)
			row[2] = &simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  durationCell,
			}
			taskCell := renderTrackedTime(trackedTime)
			row[3] = &simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  taskCell,
			}

			table.Body.Cells = append(table.Body.Cells, row)
		}

		dailySumRow := []*simpletable.Cell{
			{},
			{Align: simpletable.AlignRight, Text: "Sum:"},
			{Align: simpletable.AlignLeft, Text: "[" + dailySum.String() + "]"},
			{},
		}
		table.Body.Cells = append(table.Body.Cells, dailySumRow)

		table.Body.Cells = append(table.Body.Cells, emptyRow)

	}

	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())
}

func renderTrackedTime(te heimat.TrackEntry) string {
	return fmt.Sprintf("%s %s", te.Project.Name, te.Task.Name)
}
func renderTime(te heimat.TrackEntry) string {
	return fmt.Sprintf("%s-%s", te.Start, te.End)
}

func renderDuration(te heimat.TrackEntry) string {
	dur := calcDuration(te.Start, te.End)
	return fmt.Sprintf("%s", dur.String())
}

func calcDuration(s, e string) time.Duration {
	start, _ := time.Parse("15:04", s)
	end, _ := time.Parse("15:04", e)
	dur := end.Sub(start)
	return dur
}

type byDate []heimat.Day

func (s byDate) Len() int {
	return len(s)
}
func (s byDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byDate) Less(i, j int) bool {
	return s[i].Date < s[j].Date
}
