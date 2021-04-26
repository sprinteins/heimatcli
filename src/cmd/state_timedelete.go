package main

import (
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/api"
	"heimatcli/src/heimat/print"
	"heimatcli/src/x/log"
	"strconv"
	"strings"
	"time"

	prompt "github.com/c-bata/go-prompt"
)

// StateTimeDelete _
type StateTimeDelete struct {
	api    *api.API
	day    *heimat.Day
	date   time.Time
	cancel func()
}

var stateTimeDeleteSetTime = make(chan time.Time, 1)

// NewStateTimeDelete _
func NewStateTimeDelete(api *api.API, cancelFn func()) *StateTimeDelete {
	return &StateTimeDelete{
		api:    api,
		cancel: cancelFn,
	}
}

// Suggestions _
func (std StateTimeDelete) Suggestions(in prompt.Document) []prompt.Suggest {

	// suggest projects
	suggestions := make([]prompt.Suggest, len(std.day.TrackedTimes))
	for ti, track := range std.day.TrackedTimes {
		suggestions[ti] = std.generatetimeSuggestion(track)
	}

	return suggestions
}

func (std StateTimeDelete) generatetimeSuggestion(track heimat.TrackEntry) prompt.Suggest {
	return prompt.Suggest{
		Text:        fmt.Sprintf("%d", track.ID),
		Description: fmt.Sprintf("%s-%s %s %s", track.Start, track.End, track.Project.Name, track.Task.Name),
	}
}

// Prefix _
func (std StateTimeDelete) Prefix() string {
	if sameDay(std.date) {
		return "heimat > time delete > "
	}

	monthDay := std.date.Format("01-02")
	return fmt.Sprintf("Heimat > time delete (%s) > ", monthDay)
}

// Exe _
func (std StateTimeDelete) Exe(in string) StateKey {

	cmd := normalizeCommand(in)

	idStrings := strings.Split(cmd, " ")
	ids := make([]int, 0)
	for _, idString := range idStrings {
		id, err := strconv.Atoi(idString)
		if err != nil {
			log.Error.Printf("msg='could not convert string to int', str='%s', err='%s'", idString, err)
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return stateKeyNoChange
	}

	for _, id := range ids {
		std.api.DeleteTime(id)
	}

	day := std.api.FetchDayByDate(std.date)

	print.Day(day)

	return stateKeyNoChange
}

// Init _
func (std *StateTimeDelete) Init() {
	date := <-stateTimeDeleteSetTime
	std.date = date
	std.day = std.api.FetchDayByDate(date)
}
