package main

import (
	"fmt"
	"heimatcli/heimat"
	"heimatcli/heimat/api"
	"heimatcli/x/log"
	"regexp"
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
		start, _ := time.Parse("15:04", tt.Start)
		end, _ := time.Parse("15:04", tt.End)
		diff := end.Sub(start)
		subTime = subTime + diff

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: tt.Start},
			{Align: simpletable.AlignRight, Text: tt.End},
			{Align: simpletable.AlignLeft, Text: diff.String()},
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
