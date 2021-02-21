package api

import (
	"encoding/json"
	"heimatcli/src/heimat"
	"heimatcli/src/x/log"
	"time"
)

// FetchBalances _
func (api *API) FetchBalances(date time.Time) *heimat.Balances {
	year, _, _ := date.Date()
	url := api.urlBalances(api.UserID(), year)

	resp, _, err := api.httpGet(api.Token(), url, nil)
	if err != nil {
		log.Error.Printf("Could not fetch projects: %s\n", err.Error())
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Error.Printf("could not fetch project, HTTP Status: %d", resp.StatusCode)
		return nil
	}

	respBodyBytes := readBody(resp)
	balances := &heimat.Balances{}
	err = json.Unmarshal(respBodyBytes, balances)
	if err != nil {
		log.Error.Printf("could not unmarshal response body to balances: %s\n", err.Error())
		return nil
	}

	return balances
}
