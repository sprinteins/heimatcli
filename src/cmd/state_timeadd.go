package main

import (
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/api"
	"regexp"
	"strings"
	"time"

	prompt "github.com/c-bata/go-prompt"
)

// StateTimeAdd _
type StateTimeAdd struct {
	api      *api.API
	projects []heimat.Project
	date     time.Time
	cancel   func()
}

var stateTimeAddSetTime = make(chan time.Time, 1)

// NewStateTimeAdd _
func NewStateTimeAdd(api *api.API, cancelFn func()) *StateTimeAdd {
	return &StateTimeAdd{
		api:    api,
		cancel: cancelFn,
	}
}

// Suggestions _
func (sta StateTimeAdd) Suggestions(in prompt.Document) []prompt.Suggest {

	cmd := normalizeCommand(in.Text)

	// parse start time, suggest end time
	start, end := sta.findTimes(cmd)

	if start != "" && end != "" {
		return []prompt.Suggest{}
	}

	if start != "" {
		return timeSuggestionsWithMin(start)
	}

	// parse task, suggest start times
	project := sta.findProject(cmd)
	if project != nil {
		task := sta.findTask(*project, cmd)
		if task != nil {
			return timeSuggestions()
		}
	}

	// parse project, suggest tasks
	if project != nil {
		suggestions := make([]prompt.Suggest, len(project.Tasks))
		for ti, task := range project.Tasks {
			suggestions[ti] = prompt.Suggest{Text: task.Name}
		}
		return suggestions
	}

	// suggest projects
	suggestions := make([]prompt.Suggest, len(sta.projects))
	for pi, project := range sta.projects {
		suggestions[pi] = prompt.Suggest{Text: project.Name, Description: fmt.Sprintf("%d", (project.ID))}
	}

	return suggestions
}

// Prefix _
func (sta StateTimeAdd) Prefix() string {
	if sameDay(sta.date) {
		return "heimat > time add > "
	}

	monthDay := sta.date.Format("01.02")
	return fmt.Sprintf("Heimat > time add (%s) > ", monthDay)
}

// Exe _
func (sta StateTimeAdd) Exe(in string) StateKey {

	cmd := normalizeCommand(in)
	project := sta.findProject(cmd)
	if project == nil {
		fmt.Println("Could not find project!")
		return stateKeyNoChange
	}

	task := sta.findTask(*project, cmd)
	if task == nil {
		fmt.Println("Could not find task!")
		return stateKeyNoChange
	}

	start, end := sta.findTimes(cmd)
	if start == "" {
		fmt.Println("Could not find start time!")
		return stateKeyNoChange
	}
	if end == "" {
		fmt.Println("Could not find end time!")
		return stateKeyNoChange
	}
	if start >= end {
		fmt.Println("Start has to be earlier than end!")
		return stateKeyNoChange
	}

	note := sta.findNotes(cmd, *project, *task, start, end)
	if note == "" {
		fmt.Println("Could not found notes!")
		return stateKeyNoChange
	}

	date := sta.date
	sta.api.SendCreateTime(sta.api.UserID(), date, start, end, note, *task)

	day := sta.api.FetchDayByDate(sta.date)

	printDay(day)

	return stateKeyNoChange
}

// Init _
func (sta *StateTimeAdd) Init() {
	date := <-stateTimeAddSetTime
	sta.date = date
	sta.projects = sta.api.FetchProjects()
}

func sameDay(d time.Time) bool {
	layout := "2006-01-02"
	dStr := d.Format(layout)
	nowStr := time.Now().Format(layout)
	return dStr == nowStr
}

func (sta StateTimeAdd) findProject(cmd string) *heimat.Project {

	for _, project := range sta.projects {
		if strings.Contains(cmd, project.Name) {
			return &project
		}
	}

	return nil
}

func (sta StateTimeAdd) findTask(project heimat.Project, cmd string) *heimat.Task {

	for _, task := range project.Tasks {
		if strings.Contains(cmd, task.Name) {
			return &task
		}
	}

	return nil
}

func (sta StateTimeAdd) findTimes(cmd string) (start, end string) {
	re := regexp.MustCompile(`(\d{2}:\d{2})`)
	res := re.FindAll([]byte(cmd), -1)

	// Parse start time
	if len(res) > 0 {
		start = string(res[0])
	}

	// Parse end time
	if len(res) > 1 {
		end = string(res[1])
	}

	return start, end

}

func (sta StateTimeAdd) findNotes(
	cmd string,
	project heimat.Project,
	task heimat.Task,
	start string,
	end string,
) string {
	rest := cmd
	rest = strings.Replace(rest, project.Name, "", 1)
	rest = strings.Replace(rest, task.Name, "", 1)
	rest = strings.Replace(rest, start, "", 1)
	rest = strings.Replace(rest, end, "", 1)

	rest = strings.TrimSpace(rest)
	return rest
}
