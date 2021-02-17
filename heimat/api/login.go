package api

import (
	"encoding/json"
	"heimatcli/x/log"
	"heimatcli/x/types"
	"io/ioutil"
)

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login _
func (api API) Login(username, password string) types.Token {

	payload := loginPayload{
		Username: username,
		Password: password,
	}

	p, err := json.Marshal(payload)
	if err != nil {
		log.Error.Printf("could not marshal login payload: %s", err)
	}

	postURL := api.urlAuthentication()
	resp, _, err := httpPost("", postURL, nil, p)
	if err != nil {
		log.Error.Printf("could not make login request: %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		log.Warning.Printf(
			"msg=non-success-response url=%s code=%d body=%s resp=%+v",
			postURL,
			resp.StatusCode,
			body,
			resp,
		)
	}

	loginResp := &loginResponse{}
	err = json.Unmarshal(body, loginResp)
	if err != nil {
		log.Error.Printf("could not unmarshal login response: %s", err)
	}

	return loginResp.Token

}

// IsAuthenticated _
func (api *API) IsAuthenticated() bool {
	if api.token == "" {
		return false
	}

	return api.isTokenValid()
}

func (api *API) isTokenValid() bool {

	userID := api.UserID()
	_, err := api.FetchUserByID(userID)

	if err != nil {
		return false
	}

	return true
}
