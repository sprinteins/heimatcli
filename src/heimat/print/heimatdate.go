package print

import (
	"fmt"
	"heimatcli/src/heimat/api"
)

// HeimatDate reformats and print the heimat's date format
func HeimatDate(d string) {
	date := api.DateFromHeimatDate(d)
	dateString := date.Format("2006-01-02 (Mon)")

	fmt.Printf("\n%s\n\n", dateString)
}
