package print

import (
	"fmt"
	"time"
)

func Date(d time.Time) {
	dateString := FormatDate(d)

	fmt.Printf("\n%s\n\n", dateString)
}

func FormatDate(d time.Time) string {
	dateString := d.Format("2006-01-02 (Mon)")
	return dateString
}
