package api

import (
	"encoding/json"
	"heimatcli/heimat"
	"heimatcli/x/log"
)

// FetchUserByID _
// https://heimat-demo.sprinteins.com/api/v1/employees/42
func (api *API) FetchUserByID(userID int) (*heimat.User, error) {
	apiURL := api.urlEmployeeByID(userID)

	resp, _, err := httpGet(api.Token(), apiURL, nil)
	if err != nil {
		return nil, err
	}

	respBytes := readBody(resp)
	user := &heimat.User{}
	err = json.Unmarshal(respBytes, user)
	if err != nil {
		log.Error.Printf(
			"could not unmarshal user response: %s\n%s\n",
			err.Error(),
			string(respBytes),
		)

		return nil, err
	}

	return user, nil
}
