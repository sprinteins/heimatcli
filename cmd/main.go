package main

import (
	"heimatcli/heimat/api"

	prompt "github.com/c-bata/go-prompt"
)

func completer(in prompt.Document) []prompt.Suggest {
	s := sm.currentState.Suggestions(in)
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

var LivePrefixState struct {
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

func main() {

	heimatAPI := api.NewAPI("https://heimat.sprinteins.com/api/v1")

	sm = NewStateMachine(stateKeyHome)
	sm.AddState(stateKeyLogin, NewStateLogin(heimatAPI))
	sm.AddState(stateKeyHome, NewStateHome(heimatAPI))
	sm.AddState(stateKeyTimeAdd, NewStateTimeAdd(heimatAPI, sm.Cancel))

	if heimatAPI.IsAuthenticated() {
		sm.ChangeState(stateKeyHome)
	} else {
		sm.ChangeState(stateKeyLogin)
	}

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("live-prefix-example"),
		prompt.OptionSwitchKeyBindMode(prompt.CommonKeyBind),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.Escape,
			Fn: func(b *prompt.Buffer) {
				sm.Cancel()
			},
		}),
	)
	p.Run()
}
