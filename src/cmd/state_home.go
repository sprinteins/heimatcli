package main

import (
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/api"
	"heimatcli/src/x/log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

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

	cmd := normalizeCommand(in.Text)
	noSuggestions := []prompt.Suggest{}

	if strings.Contains(cmd, "time show") {
		return []prompt.Suggest{
			{Text: "day", Description: "Show Day"},
			{Text: "month", Description: "Show Month"},
		}
	}

	if strings.Contains(cmd, "time add") {
		return noSuggestions
	}

	if strings.Contains(cmd, "profile") {
		return noSuggestions
	}

	if strings.Contains(cmd, "time copy") {
		return noSuggestions
	}

	if strings.Contains(cmd, "time") {
		return []prompt.Suggest{
			{Text: "show", Description: "Show Tracked Time"},
			{Text: "add", Description: "Add Time"},
			{Text: "copy", Description: "Copy a day"},
		}
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
		date := dateFromCommand(cmd, "time show day")
		day := sh.api.FetchDayByDate(date)
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
		printProfile(u, b)
	}

	if strings.Contains(cmd, "time add") {
		date := dateFromCommand(cmd, "time add")
		stateTimeAddSetTime <- date
		return stateKeyTimeAdd
	}

	if strings.Contains(cmd, "time copy") {
		sourceDate, targetDate := sourceAndTargetDateFromCMD(cmd, "time copy")
		if sourceDate == nil {
			log.Error.Printf("could not determine source date")
			return stateKeyNoChange
		}

		if targetDate == nil {
			log.Error.Printf("could not determine target date")
			return stateKeyNoChange
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
		printDay(day)

		// TODO: fetch time from source date and crate on target date
		return stateKeyNoChange
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

//
// TIME ADD
//
func dateFromCommand(cmd string, strToRemove string) time.Time {
	rest := strings.Replace(cmd, strToRemove, "", 1)
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return time.Now()
	}

	return calcDateFromString(rest)
}

func isRelativeDate(d string) bool {
	relativeDate := regexp.MustCompile(`^(\+|-)\d*$`)
	return relativeDate.Match([]byte(d))
}

func calcDateFromString(dateStr string) time.Time {

	if isRelativeDate(dateStr) || dateStr == "0" {
		return calcRelativeDate(dateStr)
	}
	return calcAbsoluteDate(dateStr)

}

func calcRelativeDate(relativeDate string) time.Time {
	diff, err := strconv.Atoi(relativeDate)
	if err != nil {
		log.Error.Printf("could not parse relative date: %s\n", err)
		return time.Now()
	}
	return time.Now().AddDate(0, 0, diff)
}

func calcAbsoluteDate(absDate string) time.Time {
	day, err := strconv.Atoi(absDate)
	if err != nil {
		log.Error.Printf("could not parse into day: %s\n", err)
	}
	now := time.Now()
	year, month, _ := now.Date()
	dateStr := fmt.Sprintf("%d-%d-%d", year, month, day)
	newDate, err := time.Parse("2006-1-2", dateStr)
	if err != nil {
		log.Error.Printf("could not create new date from absolute date:%s", err)
	}
	return newDate
}

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
		s := calcDateFromString(string(dates[0]))
		t := time.Now()
		return &s, &t
	}

	if len(dates) == 2 {
		s := calcDateFromString(string(dates[0]))
		t := calcDateFromString(string(dates[1]))
		return &s, &t
	}

	return nil, nil

}
