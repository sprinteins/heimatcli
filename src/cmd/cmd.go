package main

import (
	"flag"
	"fmt"
	"heimatcli/src/heimat/api"

	prompt "github.com/c-bata/go-prompt"
)

func completer(in prompt.Document) []prompt.Suggest {
	s := sm.currentState.Suggestions(in)
	return prompt.FilterFuzzy(s, in.GetWordBeforeCursor(), true)
}

var livePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func executor(in string) {

	newStateKey := sm.currentState.Exe(in)
	if newStateKey == stateKeyNoChange {
		return
	}

	sm.ChangeState(newStateKey)

}

func changeLivePrefix() (string, bool) {
	prefix := sm.CurrentState().Prefix()
	return prefix, true
}

var sm *StateMachine

// Run start the app
func Run() {

	const defaultAPI = "https://heimat.sprinteins.com/api/v1"

	// parse flags
	apiEndpoint := flag.String("api", defaultAPI, "API Endpoint")
	versionRequest := flag.Bool("version", false, "Prints current version")

	flag.Parse()

	if *versionRequest {
		fmt.Println("v0.1.5")
		return
	}

	// Initialize Dependencies
	heimatAPI := api.NewAPI(*apiEndpoint)

	startPrompt(heimatAPI)

}
func cliLogin(api *api.API) {
	loginState := NewStateLogin(api)
	loginState.Exe("")
}

func startPrompt(heimatAPI *api.API) {

	// cancel function
	cancel := func() {
		if !heimatAPI.IsAuthenticated() {
			sm.ChangeState(stateKeyLogin)
		}
		sm.ChangeState(stateKeyHome)
	}
	cancelPrompt := func(b *prompt.Buffer) {
		cancel()
	}

	// Initialize states
	sm = NewStateMachine(stateKeyHome)
	sm.AddState(stateKeyLogin, NewStateLogin(heimatAPI))
	sm.AddState(stateKeyHome, NewStateHome(heimatAPI))
	sm.AddState(stateKeyTimeAdd, NewStateTimeAdd(heimatAPI, cancel))
	sm.AddState(stateKeyTimeDelete, NewStateTimeDelete(heimatAPI, cancel))

	// Skip Login if already authenticated
	if !heimatAPI.IsAuthenticated() {
		stateLogin := NewStateLogin(heimatAPI)
		loginSuccess := stateLogin.Login()
		for !loginSuccess {
			fmt.Println("Wrong username or password!")
			loginSuccess = stateLogin.Login()
		}
	}
	sm.ChangeState(stateKeyHome)

	// Set up go-prompt
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("Heimat"),
		prompt.OptionSwitchKeyBindMode(prompt.CommonKeyBind),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.Escape,
			Fn:  cancelPrompt,
		}),
	)
	p.Run()
}
