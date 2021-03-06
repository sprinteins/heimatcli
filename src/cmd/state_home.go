package main

import (
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/api"
	"heimatcli/src/heimat/print"
	"heimatcli/src/x/date"
	"heimatcli/src/x/log"
	"regexp"
	"strings"
	"sync"
	"time"

	prompt "github.com/c-bata/go-prompt"
)

// StateHome _
type StateHome struct {
	api               *api.API
	commands          []command
	suggestions       []suggestion
	defaultSuggestion []prompt.Suggest
}

type commandFn = func(cmd string) *StateKey
type command struct {
	key     string
	command commandFn
}

type suggestFn = func(cmd string) []prompt.Suggest
type suggestion struct {
	key     string
	suggest suggestFn
}

// NewStateHome _
func NewStateHome(api *api.API) *StateHome {

	ctrlStats := NewCtrlStats(api)

	sh := &StateHome{
		api: api,
	}
	sh.commands = []command{
		{key: "time show day", command: sh.showDay},
		{key: "time show month", command: sh.showMonth},
		{key: "time add", command: sh.changeToTimeAdd},
		{key: "time delete", command: sh.changeToTimeDelete},
		{key: "time copy", command: sh.copyTime},
		{key: "profile", command: sh.showProfile},
		{key: "stats", command: ctrlStats.ShowStats},

		{key: "logout", command: sh.logout},
	}

	noSuggestions := func(_ string) []prompt.Suggest { return []prompt.Suggest{} }
	sh.suggestions = []suggestion{
		{key: "time show day", suggest: noSuggestions},
		{key: "time show month", suggest: noSuggestions},
		{key: "time add", suggest: noSuggestions},
		{key: "time delete", suggest: noSuggestions},
		{key: "profile", suggest: noSuggestions},
		{key: "logout", suggest: noSuggestions},
		{key: "time copy", suggest: noSuggestions},

		{
			key: "time show",
			suggest: func(_ string) []prompt.Suggest {
				return []prompt.Suggest{
					{Text: "day", Description: "Show Day"},
					{Text: "month", Description: "Show Month"},
				}
			},
		},

		{
			key: "time",
			suggest: func(_ string) []prompt.Suggest {
				return []prompt.Suggest{
					{Text: "show", Description: "Show Tracked Time"},
					{Text: "add", Description: "Add Time"},
					{Text: "copy", Description: "Copy a day"},
					{Text: "delete", Description: "Delete Tracked Times"},
				}
			},
		},
		{
			key:     "stats",
			suggest: ctrlStats.Suggestions,
		},
	}

	sh.defaultSuggestion = []prompt.Suggest{
		{Text: "time", Description: "Time Tracking"},
		{Text: "profile", Description: "Show the profile and stats about the user"},
		{Text: "logout", Description: "Logout"},
		{Text: "stats", Description: "Show statistics"},
	}

	return sh
}

// Suggestions _
func (sh StateHome) Suggestions(in prompt.Document) []prompt.Suggest {

	cmd := normalizeCommand(in.Text)

	for _, suggestion := range sh.suggestions {
		if strings.Contains(cmd, suggestion.key) {
			return suggestion.suggest(cmd)
		}
	}

	return sh.defaultSuggestion

}

// Prefix _
func (sh StateHome) Prefix() string {
	return "heimat > "
}

// Exe _
func (sh StateHome) Exe(in string) StateKey {

	cmd := normalizeCommand(in)
	var newKey *StateKey

	for _, command := range sh.commands {
		if strings.Contains(cmd, command.key) {
			newKey = command.command(cmd)
			break
		}
	}

	defaultKey := stateKeyNoChange
	if newKey == nil {
		return defaultKey
	}
	return *newKey
}

func (sh StateHome) showDay(cmd string) *StateKey {
	date := date.DateFromCommand(cmd, "time show day")
	day := sh.api.FetchDayByDate(date)
	if day == nil {
		return nil
	}
	print.Day(day)
	return nil
}

func (sh StateHome) showMonth(cmd string) *StateKey {
	month := sh.api.FetchMonthByDate(time.Now())
	if month == nil {
		return nil
	}
	print.Month(month)
	return nil
}

func (sh StateHome) showProfile(cmd string) *StateKey {
	var u *heimat.User
	var b *heimat.Balances
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		u = sh.api.FetchUserByID(sh.api.UserID())
		wg.Done()
	}()

	go func() {
		b = sh.api.FetchBalances(time.Now())
		wg.Done()
	}()

	wg.Wait()
	print.Profile(u, b)
	return nil
}

func (sh StateHome) changeToTimeAdd(cmd string) *StateKey {
	date := date.DateFromCommand(cmd, "time add")
	stateTimeAddSetTime <- date
	newKey := stateKeyTimeAdd
	return &newKey
}

func (sh StateHome) changeToTimeDelete(cmd string) *StateKey {
	date := date.DateFromCommand(cmd, "time delete")
	stateTimeDeleteSetTime <- date
	newKey := stateKeyTimeDelete
	return &newKey
}

func (sh StateHome) copyTime(cmd string) *StateKey {
	sourceDate, targetDate := sourceAndTargetDateFromCMD(cmd, "time copy")
	if sourceDate == nil {
		log.Error.Printf("could not determine source date")
		return nil
	}

	if targetDate == nil {
		log.Error.Printf("could not determine target date")
		return nil
	}

	sourceDay := sh.api.FetchDayByDate(*sourceDate)

	var wg sync.WaitGroup
	wg.Add(len(sourceDay.TrackedTimes))
	for _, tt := range sourceDay.TrackedTimes {
		go func(tt heimat.TrackEntry) {
			sh.api.SendCreateTime(sh.api.UserID(), *targetDate, tt.Start, tt.End, tt.Note, tt.Task)
			wg.Done()
		}(tt)
	}
	wg.Wait()

	day := sh.api.FetchDayByDate(*targetDate)
	print.Day(day)

	return nil
}

func (sh StateHome) logout(cmd string) *StateKey {
	sh.api.Logout()
	fmt.Printf("Good bye! 👋\n")
	newKey := stateKeyLogin
	return &newKey
}

// Init noop
func (sh StateHome) Init() {}

func normalizeCommand(cmd string) string {
	singleSpacePattern := regexp.MustCompile(`\s+`)
	withSingleSpaces := singleSpacePattern.ReplaceAllString(cmd, " ")
	return strings.TrimSpace(withSingleSpaces)
}

//
// TIME ADD
//

//
// TIME COPY
//
func sourceAndTargetDateFromCMD(cmd string, baseCMD string) (source, target *time.Time) {
	rest := strings.Replace(cmd, baseCMD, "", 1)
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return nil, nil
	}

	dayDates := regexp.MustCompile(`(\+|-)?\d+`)
	dates := dayDates.FindAll([]byte(rest), 2)

	if len(dates) == 1 {
		s := date.CalcDateFromString(string(dates[0]))
		t := time.Now()
		return &s, &t
	}

	if len(dates) == 2 {
		s := date.CalcDateFromString(string(dates[0]))
		t := date.CalcDateFromString(string(dates[1]))
		return &s, &t
	}

	return nil, nil

}
