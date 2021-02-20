package main

import (
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/api"
	"heimatcli/src/x/log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	prompt "github.com/c-bata/go-prompt"
)

// StateHome _
type StateHome struct {
	api *api.API
}

// NewStateHome _
func NewStateHome(api *api.API) *StateHome {
	return &StateHome{
		api: api,
	}
}

// Suggestions _
func (sh StateHome) Suggestions(in prompt.Document) []prompt.Suggest {

	// cmd := in.GetWordBeforeCursorWithSpace()
	cmd := normalizeCommand(in.Text)

	if strings.Contains(cmd, "time show") {
		return []prompt.Suggest{
			{Text: "day", Description: "Show Day"},
			{Text: "month", Description: "Show Month"},
		}
	}

	if strings.Contains(cmd, "time add") {
		return []prompt.Suggest{}
	}

	if strings.Contains(cmd, "time") {
		return []prompt.Suggest{
			{Text: "show", Description: "Show Tracked Time"},
			{Text: "add", Description: "Add Time"},
		}
	}

	if strings.Contains(cmd, "profile") {
		return []prompt.Suggest{}
	}

	return []prompt.Suggest{
		{Text: "time", Description: "Time Tracking"},
		{Text: "profile", Description: "Show the profile and stats about the user"},
	}
}

// Prefix _
func (sh StateHome) Prefix() string {
	return "heimat > "
}

// Exe _
func (sh StateHome) Exe(in string) StateKey {

	cmd := normalizeCommand(in)

	if strings.Contains(cmd, "time show day") {

		day := sh.api.FetchDayByDate(time.Now())
		if day == nil {
			return stateKeyNoChange
		}
		printDay(day)

	}

	if strings.Contains(cmd, "time show month") {
		month := sh.api.FetchMonthByDate(time.Now())
		if month == nil {
			return stateKeyNoChange
		}
		printMonth(month)
	}

	if strings.Contains(cmd, "profile") {
		u := sh.api.FetchUserByID(sh.api.UserID())
		printUser(u)
	}

	if strings.Contains(cmd, "time add") {
		return stateKeyTimeAdd
	}

	return stateKeyNoChange
}

// Init noop
func (sh StateHome) Init() {}

func normalizeCommand(cmd string) string {
	singleSpacePattern := regexp.MustCompile(`\s+`)
	withSingleSpaces := singleSpacePattern.ReplaceAllString(cmd, " ")
	return strings.TrimSpace(withSingleSpaces)
}

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
		// start, _ := time.Parse("15:04", tt.Start)
		// end, _ := time.Parse("15:04", tt.End)
		// diff := end.Sub(start)
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

func printHeimatDate(d string) {
	date, err := time.Parse("2006-01-02", d)
	if err != nil {
		log.Error.Printf("could not parse heimat date: %s", err)
	}
	dateString := date.Format("2006-01-02 (Mon)")

	fmt.Printf("\n%s\n\n", dateString)
}

func printUser(u *heimat.User) {
	fmt.Printf("%#v\n", u)
}

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
