package date

import (
	"fmt"
	"heimatcli/src/x/log"
	"strconv"
	"strings"
	"time"
)

func FirstLastOfMonth(date time.Time) (firstDay, lastDay time.Time) {
	now := date
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstDay = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastDay = firstDay.AddDate(0, 1, -1)

	return firstDay, lastDay
}

func MonthFromCommand(cmd string, strToRemove string) time.Time {
	rest := strings.Replace(cmd, strToRemove, "", 1)
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return time.Now()
	}

	return CalcMonthFromString(rest)
}

func CalcMonthFromString(dateStr string) time.Time {

	if IsRelativeDate(dateStr) || dateStr == "0" {
		return CalcRelativeMonth(dateStr)
	}
	return CalcAbsoluteMonth(dateStr)

}

func CalcRelativeMonth(relativeDate string) time.Time {
	diff, err := strconv.Atoi(relativeDate)
	if err != nil {
		log.Error.Printf("could not parse relative date: %s\n", err)
		return time.Now()
	}
	return time.Now().AddDate(0, diff, 1)
}

func CalcAbsoluteMonth(absDate string) time.Time {
	month, err := strconv.Atoi(absDate)
	if err != nil {
		log.Error.Printf("could not parse into month: %s\n", err)
	}
	now := time.Now()
	year, _, _ := now.Date()
	dateStr := fmt.Sprintf("%d-%d-%d", year, month, 1)
	newDate, err := time.Parse("2006-1-2", dateStr)
	if err != nil {
		log.Error.Printf("could not create new date from absolute month:%s", err)
	}
	return newDate
}
