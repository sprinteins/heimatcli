package api

import (
	"encoding/json"
	"fmt"
	"heimatcli/heimat"
	"heimatcli/x/log"
	"time"
)

// FetchDayByDate _
// https://heimat-demo.sprinteins.com/api/v1/employees/42/trackedtimes?start=2020-01-01&end=2020-01-31
func (api *API) FetchDayByDate(date time.Time) *heimat.Day {
	// userID := api.UserID()
	// apiURL := fmt.Sprintf("/employees/%d/trackedtimes?start=%s&end=%s", userID, date, date)

	// url := api.baseURL + apiURL
	url := api.urlDayByDate(api.UserID())

	queries := []Query{
		{key: "start", value: NewHeimatDate(date)},
		{key: "end", value: NewHeimatDate(date)},
	}

	resp, _, err := httpGet(api.Token(), url, queries)

	if err != nil {
		log.Error.Printf("Could not fetch projects: %s\n", err.Error())
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Error.Printf("could not fetch project, HTTP Status: %d", resp.StatusCode)
		return nil
	}

	respBodyBytes := readBody(resp)
	trackedTimes := &trackedTimeResponse{}
	err = json.Unmarshal(respBodyBytes, trackedTimes)
	if err != nil {
		log.Error.Printf("could not unmarshal response body: %s\n", err.Error())
		return nil
	}

	day := findDayByDate(trackedTimes.TrackedTimesDate, NewHeimatDate(date))

	return day

}

func findDayByDate(days []heimat.Day, date string) *heimat.Day {
	for _, day := range days {
		if day.Date == date {
			return &day
		}
	}

	return nil

}

func firstLastOfMonth() (firstDay string, lastDay string) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	fmt.Println(firstOfMonth)
	fmt.Println(lastOfMonth)

	firstDay = firstOfMonth.Format("2006-01-02")
	lastDay = lastOfMonth.Format("2006-01-02")

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
