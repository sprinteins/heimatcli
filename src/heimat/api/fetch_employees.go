package api

import (
	"encoding/json"
	"heimatcli/src/heimat"
	"heimatcli/src/x/log"
)

// FetchEmployees _
// https://heimat-demo.sprinteins.com/api/v1/employees/
func (api *API) FetchEmployees() []heimat.User {
	apiURL := api.urlEmployeesAll()

	resp, _, err := api.httpGet(api.Token(), apiURL, nil)
	if err != nil {
		log.Error.Printf("could not fetch user: %s", err)
		return nil
	}

	respBytes := readBody(resp)
	empResp := &allEmployeeResponse{}
	err = json.Unmarshal(respBytes, &empResp)
	if err != nil {
		log.Error.Printf(
			"could not unmarshal users response: err=%s\nurl=%s\nresp=%s\n",
			err.Error(),
			apiURL,
			string(respBytes),
		)

		return nil
	}

	return empResp.Employees
}

type allEmployeeResponse struct {
	Employees []heimat.User `json:"employees"`
}
