package api

import (
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/x/log"
	"net/url"
	"time"
)

const (
	urlAuthentication     = "/authentication"
	urlOidcAuthentication = "/authentication/oidc"
	urlEmployees          = "/employees"
	urlTrackedTime        = "/trackedtimes"
)

func (api API) urlAuthentication() string {

	apiURL := urlAuthentication

	fullURL := fmt.Sprintf("%s%s", api.baseURL.String(), apiURL)
	return fullURL
}

func (api API) urlOpenId() string {

	apiURL := urlOidcAuthentication

	fullURL := fmt.Sprintf("%s%s", api.baseURL.String(), apiURL)
	return fullURL
}

func (api API) urlEmployeeByID(id int) string {

	apiURL := fmt.Sprintf("%s/%d", urlEmployees, id)

	fullURL := fmt.Sprintf("%s%s", api.baseURL.String(), apiURL)
	return fullURL
}

func (api API) urlEmployeesAll() string {

	fullURL := fmt.Sprintf("%s%s", api.baseURL.String(), urlEmployees)
	return fullURL
}

func (api API) urlDayByDate(userID int) string {

	apiURL := fmt.Sprintf("/employees/%d/trackedtimes", userID)
	fullURL := fmt.Sprintf("%s/%s", api.baseURL.String(), apiURL)
	return fullURL
}

func (api API) urlProjects(userID int, date time.Time) string {

	apiPath := fmt.Sprintf("/employees/%d/projects", userID)
	fullPath := fmt.Sprintf("%s%s", api.baseURL.String(), apiPath)
	fullURL, err := url.Parse(fullPath)
	if err != nil {
		log.Error.Printf("could not parse url: %s", err)
	}

	q := fullURL.Query()
	q.Add("date", heimat.NewHeimatDate(date))
	fullURL.RawQuery = q.Encode()

	return fullURL.String()

}

// urlTrackedTimeCreate _
func (api API) urlTrackedTimeCreate() string {
	apiURL := urlTrackedTime
	baseURL := api.baseURL

	fullURL := fmt.Sprintf("%s%s", baseURL.String(), apiURL)

	return fullURL
}

func (api API) urlBalances(userID int, year int) string {
	return fmt.Sprintf("%s/employees/%d/absence/balances/%d", api.baseURL.String(), userID, year)
}

func (api API) urlTrackedTimeDelete(timeID int) string {
	apiURL := urlTrackedTime
	baseURL := api.baseURL

	fullURL := fmt.Sprintf("%s%s/%d", baseURL.String(), apiURL, timeID)
	return fullURL
}
