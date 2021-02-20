package api

import (
	"encoding/json"
	"heimatcli/src/heimat"
	"heimatcli/src/x/log"
	"time"
)

// FetchMonthByDate _
func (api *API) FetchMonthByDate(date time.Time) *heimat.Month {
	start, end := firstLastOfMonth(date)

	return api.fetchDaysByDates(start, end)
}

// FetchDayByDate _
// https://heimat-demo.sprinteins.com/api/v1/employees/42/trackedtimes?start=2020-01-01&end=2020-01-31
func (api *API) FetchDayByDate(date time.Time) *heimat.Day {

	days := api.fetchDaysByDates(date, date)
	day := findDayByDate(days.TrackedTimesDate, NewHeimatDate(date))

	return day

}

func (api *API) fetchDaysByDates(start, end time.Time) *heimat.Month {
	url := api.urlDayByDate(api.UserID())

	queries := []Query{
		{key: "start", value: NewHeimatDate(start)},
		{key: "end", value: NewHeimatDate(end)},
	}

	resp, _, err := api.httpGet(api.Token(), url, queries)

	if err != nil {
		log.Error.Printf("Could not fetch projects: %s\n", err.Error())
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Error.Printf("could not fetch project, HTTP Status: %d", resp.StatusCode)
		return nil
	}

	respBodyBytes := readBody(resp)
	month := &heimat.Month{}
	err = json.Unmarshal(respBodyBytes, month)
	if err != nil {
		log.Error.Printf("could not unmarshal response body: %s\n", err.Error())
		return nil
	}

	return month
}

func findDayByDate(days []heimat.Day, date string) *heimat.Day {
	for _, day := range days {
		if day.Date == date {
			return &day
		}
	}

	return nil

}

func firstLastOfMonth(date time.Time) (firstDay, lastDay time.Time) {
	now := date
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstDay = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastDay = firstDay.AddDate(0, 1, -1)

	return firstDay, lastDay
}

type trackedTimeResponse struct {
	TrackedTimesDate []heimat.Day `json:"trackedTimesDate"`
}

// HeimatDate _
type HeimatDate = string

// NewHeimatDate _
func NewHeimatDate(date time.Time) HeimatDate {
	return HeimatDate(date.Format("2006-01-02"))
}

// DateFromHeimatDate _
func DateFromHeimatDate(heimatDate string) time.Time {
	d, err := time.Parse("2006-01-02", heimatDate)
	if err != nil {
		log.Error.Printf("could not parse heimat date to date: %s", err)
		return time.Time{}
	}
	return d
}
