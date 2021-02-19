package main

import (
	"fmt"

	prompt "github.com/c-bata/go-prompt"
)

const (
	timeLimitStart = "05:59"
	timeLimitEnd   = "22:01"
)

type minuteStep int

const (
	minuteStep0  minuteStep = 0
	minuteStep15 minuteStep = 15
	minuteStep30 minuteStep = 30
	minuteStep45 minuteStep = 45
)

var minuteSteps = []minuteStep{
	minuteStep0,
	minuteStep15,
	minuteStep30,
	minuteStep45,
}

func timeSuggestions() []prompt.Suggest {

	return timeSuggestionsWithMinMax(timeLimitStart, timeLimitEnd)
}

func timeSuggestionsWithMin(min string) []prompt.Suggest {
	return timeSuggestionsWithMinMax(min, timeLimitEnd)
}

func timeSuggestionsWithMinMax(min, max string) []prompt.Suggest {
	suggestions := make([]prompt.Suggest, 0)

	// TODO: optimize start and end
	for hi := 0; hi <= 23; hi++ {
		for _, m := range minuteSteps {
			s := fmt.Sprintf("%02d:%02d", hi, m)
			if s > min && s < max {
				suggestions = append(suggestions, prompt.Suggest{Text: s})
			}
		}
	}

	return suggestions
}
