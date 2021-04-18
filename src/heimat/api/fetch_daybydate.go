package api

import (
	"encoding/json"
	"heimatcli/src/heimat"
	"heimatcli/src/x/date"
	"heimatcli/src/x/log"
	"time"
)

// FetchMonthByDate _
func (api *API) FetchMonthByDate(d time.Time) *heimat.Month {
	start, end := date.FirstLastOfMonth(d)

	return api.FetchDaysByDates(start, end)
}

// FetchDayByDate _
// https://heimat-demo.sprinteins.com/api/v1/employees/42/trackedtimes?start=2020-01-01&end=2020-01-31
func (api *API) FetchDayByDate(date time.Time) *heimat.Day {

	days := api.FetchDaysByDates(date, date)
	if days == nil {
		return nil
	}
	day := findDayByDate(days.TrackedTimesDate, heimat.NewHeimatDate(date))

	return day

}

func (api *API) FetchDaysByDates(start, end time.Time) *heimat.Month {
	url := api.urlDayByDate(api.UserID())

	queries := []Query{
		{key: "start", value: heimat.NewHeimatDate(start)},
		{key: "end", value: heimat.NewHeimatDate(end)},
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

type trackedTimeResponse struct {
	TrackedTimesDate []heimat.Day `json:"trackedTimesDate"`
}
