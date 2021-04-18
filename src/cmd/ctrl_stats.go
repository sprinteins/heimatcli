package main

import (
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/api"
	"heimatcli/src/heimat/calc"
	"heimatcli/src/heimat/print"
	"heimatcli/src/x/date"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
)

type CtrlStats struct {
	groupTypeSuggestions []prompt.Suggest
	timeTypeSuggestions  []prompt.Suggest
	defaultGroupType     string
	defaultTimeType      string
	api                  *api.API
}

func NewCtrlStats(api *api.API) CtrlStats {
	return CtrlStats{
		api: api,
		groupTypeSuggestions: []prompt.Suggest{
			{Text: "project", Description: "[Default] Statistics by project"},
			{Text: "task", Description: "Statistics by task"},
			{Text: "desc", Description: "Statistics by description"},
		},
		timeTypeSuggestions: []prompt.Suggest{
			// {Text: "month", Description: "[Default] Gather statistics by months"},
			// {Text: "year", Description: "Gather statistics by years"},
		},
		defaultGroupType: "project",
		defaultTimeType:  "month",
	}
}

//
// Suggestions
//

func (cts CtrlStats) Suggestions(cmd string) []prompt.Suggest {
	remaningSuggestions := make([]prompt.Suggest, 0)

	if !cts.hasSomeSuggestions(cmd, cts.timeTypeSuggestions) {
		remaningSuggestions = append(remaningSuggestions, cts.timeTypeSuggestions...)
	}

	if !cts.hasSomeSuggestions(cmd, cts.groupTypeSuggestions) {
		remaningSuggestions = append(remaningSuggestions, cts.groupTypeSuggestions...)
	}

	return remaningSuggestions
}

func (cts CtrlStats) hasSomeSuggestions(cmd string, suggestions []prompt.Suggest) bool {
	for _, suggestion := range suggestions {
		if strings.Contains(cmd, suggestion.Text) {
			return true
		}
	}

	return false
}

//
// Exec
//

func (cts CtrlStats) ShowStats(cmd string) *StateKey {
	groupType := cts.extractGroupType(cmd)
	groupType = toDefaultString(groupType, cts.defaultGroupType)
	cmd = strings.Replace(cmd, groupType, "", -1)

	timeType := cts.extractTimeType(cmd)
	timeType = toDefaultString(timeType, cts.defaultTimeType)
	cmd = strings.Replace(cmd, timeType, "", -1)

	cmd = normalizeCommand(cmd)

	start, end := cts.extractDateRangeByMonth(cmd, "stats")

	// Execute Group Type
	if groupType == "project" {
		projectStats := cts.calcProjectStats(start, end)
		print.ProjectStats(projectStats, start, end)
	}
	if groupType == "task" {
		taskStats := cts.calcTaskStats(start, end)
		print.TaskStats(taskStats, start, end)
	}

	if groupType == "desc" {
		taskStats := cts.calcDescStats(start, end)
		print.DescStats(taskStats, start, end)
	}

	return nil

}

func (cts CtrlStats) extractGroupType(cmd string) string {
	for _, suggestion := range cts.groupTypeSuggestions {
		if strings.Contains(cmd, suggestion.Text) {
			return suggestion.Text
		}
	}

	return ""
}

func (cts CtrlStats) extractTimeType(cmd string) string {
	for _, suggestion := range cts.timeTypeSuggestions {
		if strings.Contains(cmd, suggestion.Text) {
			return suggestion.Text
		}
	}

	return ""
}

func toDefaultString(original string, defaultValue string) string {
	if original == "" {
		return defaultValue
	}

	return original
}

//
// Project Stats
//

func (cts CtrlStats) calcProjectStats(start, end time.Time) heimat.ProjectStats {
	month := cts.api.FetchDaysByDates(start, end)

	buckets := make(heimat.ProjectStats, 0)
	var totalTimeSpent time.Duration = 0

	// Abolute values
	for _, day := range month.TrackedTimesDate {
		for _, tt := range day.TrackedTimes {
			// Probably vacation records do not have projects
			if tt.Project.Name == "" {
				continue
			}

			stat, ok := buckets[tt.Project.Name]
			if !ok {
				stat = heimat.ProjectStatsTimeSpent{
					Absolute: 0,
					Relative: 0,
				}
				buckets[tt.Project.Name] = stat

			}

			dur := calc.Duration(tt.Start, tt.End)
			stat.Absolute += dur
			totalTimeSpent += dur
			buckets[tt.Project.Name] = stat
		}
	}

	// Relative Values
	for key, bucket := range buckets {
		bucket.Relative = (float32(bucket.Absolute) / float32(totalTimeSpent)) * 100
		buckets[key] = bucket
	}

	return buckets
}

//
// Task Stats
//

func (cts CtrlStats) calcTaskStats(start, end time.Time) heimat.TaskStats {
	month := cts.api.FetchDaysByDates(start, end)

	buckets := make(heimat.TaskStats, 0)
	var totalTimeSpent time.Duration = 0

	// Abolute values
	for _, day := range month.TrackedTimesDate {
		for _, tt := range day.TrackedTimes {
			// Probably vacation records do not have projects
			key := tt.Task.Name
			if key == "" {
				continue
			}

			stat, ok := buckets[key]
			if !ok {
				stat = heimat.TaskStatsTimeSpent{
					ProjectName: tt.Project.Name,
					Absolute:    0,
					Relative:    0,
				}
				buckets[key] = stat

			}

			dur := calc.Duration(tt.Start, tt.End)
			stat.Absolute += dur
			totalTimeSpent += dur
			buckets[key] = stat
		}
	}

	// Relative Values
	for key, bucket := range buckets {
		bucket.Relative = (float32(bucket.Absolute) / float32(totalTimeSpent)) * 100
		buckets[key] = bucket
	}

	return buckets
}

//
// Desc Stats
//

func (cts CtrlStats) calcDescStats(start, end time.Time) heimat.DescStats {
	month := cts.api.FetchDaysByDates(start, end)

	buckets := make(heimat.DescStats, 0)
	var totalTimeSpent time.Duration = 0

	// Abolute values
	for _, day := range month.TrackedTimesDate {
		for _, tt := range day.TrackedTimes {
			// Probably vacation records do not have projects
			key := tt.Note
			if key == "" {
				continue
			}

			stat, ok := buckets[key]
			if !ok {
				stat = heimat.DescStatsTimeSpent{
					ProjectName: tt.Project.Name,
					TaskName:    tt.Task.Name,
					Absolute:    0,
					Relative:    0,
				}
				buckets[key] = stat

			}

			dur := calc.Duration(tt.Start, tt.End)
			stat.Absolute += dur
			totalTimeSpent += dur
			buckets[key] = stat
		}
	}

	// Relative Values
	for key, bucket := range buckets {
		bucket.Relative = (float32(bucket.Absolute) / float32(totalTimeSpent)) * 100
		buckets[key] = bucket
	}

	return buckets
}

//
// Date Range Extraction
//

func (cts CtrlStats) extractDateRangeByMonth(cmd string, strToRemove string) (start, end time.Time) {
	startString, endString := date.StartAndEndMonthStringFromCMD(cmd, strToRemove)

	startMonth := date.MonthFromCommand(startString, "")

	endMonth := startMonth
	if endString != "" {
		endMonth = date.MonthFromCommand(endString, "")
	}

	start, _ = date.FirstLastOfMonth(startMonth)
	_, end = date.FirstLastOfMonth(endMonth)

	return start, end
}
