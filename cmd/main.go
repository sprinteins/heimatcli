package main

import (
	"heimatcli/heimat/api"

	prompt "github.com/c-bata/go-prompt"
)

func completer(in prompt.Document) []prompt.Suggest {
	// fmt.Printf("completer in: %s\n", in.Text)
	s := sm.currentState.Suggestions(in)
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func executor(in string) {
	sm.currentState.Exe(in)
	LivePrefixState.LivePrefix = sm.CurrentState().Prefix() + "> "
	LivePrefixState.IsEnable = true

}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}

var sm *StateMachine

func main() {

	heimatAPI := api.NewAPI("https://heimat.sprinteins.com/api/v1")

	sm = NewStateMachine()
	sm.AddState("login", NewStateLogin(heimatAPI))

	sm.ChangeState("login")

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("live-prefix-example"),
	)
	p.Run()
}
