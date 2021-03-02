package heimat

import (
	"heimatcli/src/x/log"
	"time"
)

// HeimatDate is the text representation of date in Heimat
type HeimatDate = string

// NewHeimatDate create a text representation of a date in the Heimat-Date format
func NewHeimatDate(date time.Time) HeimatDate {
	return HeimatDate(date.Format("2006-01-02"))
}

// ParseHeimatDate parses a heimat-date format
func ParseHeimatDate(heimatDate string) time.Time {
	d, err := time.Parse("2006-01-02", heimatDate)
	if err != nil {
		log.Error.Printf("could not parse heimat date to date: %s", err)
		return time.Time{}
	}
	return d
}
