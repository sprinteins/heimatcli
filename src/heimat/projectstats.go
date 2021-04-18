package heimat

import (
	"sort"
	"time"
)

// ProjectStats represents gather statstics about time spent on projects
type ProjectStats map[string]ProjectStatsTimeSpent

type ProjectStatsTimeSpent struct {
	Absolute time.Duration
	Relative float32
}

// Sorting

type ProjectStatsPair struct {
	ProjectName      string
	ProjectTimeSpent ProjectStatsTimeSpent
}

func (ps ProjectStats) Sorted() []ProjectStatsPair {
	pairs := make([]ProjectStatsPair, 0)
	for key, value := range ps {
		pairs = append(pairs, ProjectStatsPair{ProjectName: key, ProjectTimeSpent: value})
	}
	sort.Sort(byProjectSpentAbsolute(pairs))
	return pairs
}

type byProjectSpentAbsolute []ProjectStatsPair

func (s byProjectSpentAbsolute) Len() int {
	return len(s)
}
func (s byProjectSpentAbsolute) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byProjectSpentAbsolute) Less(i, j int) bool {
	return s[i].ProjectTimeSpent.Absolute > s[j].ProjectTimeSpent.Absolute
}
