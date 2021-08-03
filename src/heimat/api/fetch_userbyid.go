package api

import (
	"encoding/json"
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/x/log"
)

// FetchUserByID _
// https://heimat-demo.sprinteins.com/api/v1/employees/42
func (api *API) FetchUserByID(userID int) *heimat.User {
	user, err := api.fetchUserByID(userID)
	if err != nil {
		log.Error.Print(err)
	}

	return user
}

func (api *API) fetchUserByID(userID int) (*heimat.User, error) {
	apiURL := api.urlEmployeeByID(userID)

	resp, _, err := api.httpGet(api.Token(), apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user: %s", err)

	}

	respBytes := readBody(resp)
	user := &heimat.User{}
	err = json.Unmarshal(respBytes, user)
	if err != nil {
		return nil, fmt.Errorf(
			"could not unmarshal user response: %s\n%s",
			err.Error(),
			string(respBytes),
		)
	}

	return user, nil
}
