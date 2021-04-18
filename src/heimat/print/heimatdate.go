package print

import (
	"heimatcli/src/heimat"
)

// HeimatDate reformats and print the heimat's date format
func HeimatDate(d string) {
	date := heimat.ParseHeimatDate(d)
	Date(date)
}
