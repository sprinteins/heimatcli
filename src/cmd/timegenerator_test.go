package main

import (
	"heimatcli/src/x/assert"
	"testing"

	"github.com/c-bata/go-prompt"
)

func Test_Time_Generator(t *testing.T) {
	type testCase struct {
		desc        string
		min         string
		max         string
		suggestions []prompt.Suggest
	}

	var tests = []testCase{
		{
			desc: "min and max are not included",
			min:  "05:59",
			max:  "07:01",
			suggestions: []prompt.Suggest{
				{Text: "06:00"},
				{Text: "06:15"},
				{Text: "06:30"},
				{Text: "06:45"},
				{Text: "07:00"},
			},
		},
	}

	test := func(t *testing.T, tc testCase) {

		//
		// Setup
		//

		//
		// Action
		//
		suggestions := timeSuggestionsWithMinMax(tc.min, tc.max)

		//
		// Validation
		//
		assert.Equals(t, suggestions, tc.suggestions, "not the same suggestions")

		// TODO: Test code comes here

	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) { test(t, tc) })
	}

}
