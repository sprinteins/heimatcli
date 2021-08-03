package main

import (
	"bufio"
	"fmt"
	"heimatcli/src/heimat/api"
	"heimatcli/src/x/log"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"golang.org/x/crypto/ssh/terminal"
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

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Error.Println(err)
		return false
	}

	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Error.Println(err)
		return false
	}
	fmt.Println()

	token, err := sl.api.Login(
		strings.TrimSpace(username),
		strings.TrimSpace(string(password)),
	)
	if err != nil {
		return false
	}
	sl.api.SetToken(token)

	return true
}

// Init noop
func (sl StateLogin) Init() {}
