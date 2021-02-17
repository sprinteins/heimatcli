package main

import (
	"bufio"
	"fmt"
	"heimatcli/heimat/api"
	"heimatcli/x/log"
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
	return "Login"
}

// Exe _
func (sl StateLogin) Exe(in string) {

	if sl.api.IsAuthenticated() {
		log.Debug.Println("already logged in")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Error.Println(err)
	}

	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Error.Println(err)
	}
	fmt.Println()

	token := sl.api.Login(
		strings.TrimSpace(username),
		strings.TrimSpace(string(password)),
	)
	sl.api.SetToken(token)
	log.Debug.Printf("token: %s", token)

}
