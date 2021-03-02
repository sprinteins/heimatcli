package print

import (
	"fmt"
	"heimatcli/src/heimat"
)

// HeimatDate reformats and print the heimat's date format
func HeimatDate(d string) {
	date := heimat.ParseHeimatDate(d)
	dateString := date.Format("2006-01-02 (Mon)")

	fmt.Printf("\n%s\n\n", dateString)
}
