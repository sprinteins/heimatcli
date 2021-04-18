package heimat

import (
	"sort"
	"time"
)

// ProjectStats represents gather statstics about time spent on projects
type TaskStats map[string]TaskStatsTimeSpent

type TaskStatsTimeSpent struct {
	ProjectName string
	Absolute    time.Duration
	Relative    float32
}

type TaskStatsPair struct {
	TaskName      string
	TaskTimeSpent TaskStatsTimeSpent
}

func (ts TaskStats) Sorted() []TaskStatsPair {
	pairs := make([]TaskStatsPair, 0)
	for key, value := range ts {
		pairs = append(pairs, TaskStatsPair{TaskName: key, TaskTimeSpent: value})
	}
	sort.Sort(byTaskSpentAbsolute(pairs))
	return pairs
}

type byTaskSpentAbsolute []TaskStatsPair

func (s byTaskSpentAbsolute) Len() int {
	return len(s)
}
func (s byTaskSpentAbsolute) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byTaskSpentAbsolute) Less(i, j int) bool {
	return s[i].TaskTimeSpent.Absolute > s[j].TaskTimeSpent.Absolute
}
