package date

import (
	"fmt"
	"heimatcli/src/x/log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ---------
// -- DAY --
// ---------

func DateFromCommand(cmd string, strToRemove string) time.Time {
	rest := strings.Replace(cmd, strToRemove, "", 1)
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return time.Now()
	}

	return CalcDateFromString(rest)
}

func CalcDateFromString(dateStr string) time.Time {

	if IsRelativeDate(dateStr) || dateStr == "0" {
		return CalcRelativeDate(dateStr)
	}
	return CalcAbsoluteDate(dateStr)

}

func IsRelativeDate(d string) bool {
	relativeDate := regexp.MustCompile(`^(\+|-)\d*$`)
	return relativeDate.Match([]byte(d))
}

func CalcRelativeDate(relativeDate string) time.Time {
	diff, err := strconv.Atoi(relativeDate)
	if err != nil {
		log.Error.Printf("could not parse relative date: %s\n", err)
		return time.Now()
	}
	return time.Now().AddDate(0, 0, diff)
}

func CalcAbsoluteDate(absDate string) time.Time {
	day, err := strconv.Atoi(absDate)
	if err != nil {
		log.Error.Printf("could not parse into day: %s\n", err)
	}
	now := time.Now()
	year, month, _ := now.Date()
	dateStr := fmt.Sprintf("%d-%d-%d", year, month, day)
	newDate, err := time.Parse("2006-1-2", dateStr)
	if err != nil {
		log.Error.Printf("could not create new date from absolute date:%s", err)
	}
	return newDate
}

func StartAndEndMonthStringFromCMD(cmd string, baseCMD string) (source, target string) {
	rest := strings.Replace(cmd, baseCMD, "", 1)
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return "", ""
	}

	dayDates := regexp.MustCompile(`(\+|-)?\d+`)
	dates := dayDates.FindAll([]byte(rest), 2)

	if len(dates) == 1 {
		s := string(dates[0])
		t := string(dates[0])
		return s, t
	}

	if len(dates) == 2 {
		s := string(dates[0])
		t := string(dates[1])
		return s, t
	}

	return "", ""

}
