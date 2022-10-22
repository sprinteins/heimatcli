package api

import (
	"heimatcli/src/x/log"
	"net/url"
)

// API _
type API struct {
	baseURL  url.URL
	clientId string
	token    string
}

// NewAPI _
func NewAPI(baseURL string, clntId string) *API {

	base, err := url.Parse(baseURL)
	if err != nil {
		log.Error.Printf("could not parse url: %s", err)
		return nil
	}

	heimatAPI := &API{
		baseURL:  *base,
		clientId: clntId,
	}
	heimatAPI.loadToken()

	return heimatAPI
}

// SetToken _
func (api *API) SetToken(token string) {
	api.token = token
	saveToken(api.token)
}

func (api *API) loadToken() {
	token := readToken()
	api.token = token
}

// Token _
func (api *API) Token() string {
	return api.token
}

// UserID returns the logged in user's id
func (api *API) UserID() int {
	userID := ExtractEmployeeID(api.token)

	return userID
}

// Client ID
func (api *API) ClientId() string {
	return api.clientId
}
