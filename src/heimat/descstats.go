package heimat

import (
	"sort"
	"time"
)

// DescStats represents gather statstics about time spent by description
type DescStats map[string]DescStatsTimeSpent

type DescStatsTimeSpent struct {
	ProjectName string
	TaskName    string
	Absolute    time.Duration
	Relative    float32
}

type DescStatsPair struct {
	Desc          string
	DescTimeSpent DescStatsTimeSpent
}

func (ds DescStats) Sorted() []DescStatsPair {
	pairs := make([]DescStatsPair, 0)
	for key, value := range ds {
		pairs = append(pairs, DescStatsPair{Desc: key, DescTimeSpent: value})
	}
	sort.Sort(byDescSpentAbsolute(pairs))
	return pairs
}

type byDescSpentAbsolute []DescStatsPair

func (s byDescSpentAbsolute) Len() int {
	return len(s)
}
func (s byDescSpentAbsolute) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byDescSpentAbsolute) Less(i, j int) bool {
	return s[i].DescTimeSpent.Absolute > s[j].DescTimeSpent.Absolute
}
