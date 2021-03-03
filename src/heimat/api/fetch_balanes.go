package api

import (
	"encoding/json"
	"heimatcli/src/heimat"
	"heimatcli/src/x/log"
	"time"
)

// FetchBalances _
func (api *API) FetchBalances(date time.Time) *heimat.Balances {
	return api.FetchBalancesByUser(api.UserID(), date)
}

func (api *API) FetchBalancesByUser(userID int, date time.Time) *heimat.Balances {
	year, _, _ := date.Date()
	url := api.urlBalances(userID, year)

	resp, _, err := api.httpGet(api.Token(), url, nil)
	if err != nil {
		log.Error.Printf("Could not fetch projects: %s\n", err.Error())
		return nil
	}
	respBodyBytes := readBody(resp)

	if resp.StatusCode >= 300 {
		log.Error.Printf(
			"msg='could not fetch balances' http_status=%d url='%s' resp='%s'",
			resp.StatusCode,
			url,
			string(respBodyBytes),
		)
		return nil
	}

	balances := &heimat.Balances{}
	err = json.Unmarshal(respBodyBytes, balances)
	if err != nil {
		log.Error.Printf("could not unmarshal response body to balances: %s\n", err.Error())
		return nil
	}

	return balances
}
