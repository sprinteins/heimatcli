package calc

import "time"

// Duration calculates the duration between two time point
// represented by the `15:04` format
//
// For examples see `calc_test.go`
func Duration(s, e string) time.Duration {
	start, _ := time.Parse("15:04", s)
	end, _ := time.Parse("15:04", e)
	dur := end.Sub(start)
	return dur
}
