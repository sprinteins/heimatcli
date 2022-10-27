package main

import (
	"context"
	"fmt"
	"heimatcli/src/heimat/api"
	"heimatcli/src/x/log"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	prompt "github.com/c-bata/go-prompt"
)

// StateLogin _
type StateLogin struct {
	api *api.API
}

// NewStateLogin _
func NewStateLogin(api *api.API) *StateLogin {
	return &StateLogin{
		api: api,
	}
}

// Suggestions _
func (sl StateLogin) Suggestions(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "login", Description: "Login into Heimat"},
	}
}

// Prefix _
func (sl StateLogin) Prefix() string {
	return "Login > "
}

// Exe _
func (sl StateLogin) Exe(in string) StateKey {

	ok := sl.Login()
	if !ok {
		return stateKeyNoChange
	}

	return stateKeyHome

}

func (sl StateLogin) Login() bool {
	if sl.api.IsAuthenticated() {
		return true
	}

	fmt.Println("No session found. Opening browser to login.")

	azureAdAuthority := "https://login.microsoftonline.com/331e8350-a57c-43c3-9a37-d76cf8000f52"
	publClientApp, publClientAppErr := public.New(sl.api.ClientId(), public.WithAuthority(azureAdAuthority))

	if publClientAppErr != nil {
		panic(publClientAppErr)
	}

	authResult, authErr := publClientApp.AcquireTokenInteractive(context.TODO(), nil)

	if authErr != nil {
		log.Error.Fatalf("An error occurred during OIDC authorization: %s", authErr)
	}

	token, tokenErr := sl.api.Login(
		strings.TrimSpace(authResult.IDToken.RawToken),
	)

	if tokenErr != nil {
		log.Warning.Fatalf("An error occurred during Exchange of tokens: %s", authErr)
		return false
	}
	sl.api.SetToken(token)

	fmt.Println("Login was successful.")

	return true
}

// Init noop
func (sl StateLogin) Init() {}
