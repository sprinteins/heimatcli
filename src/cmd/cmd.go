package main

import (
	"flag"
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/api"
	"heimatcli/src/heimat/print"
	"sync"
	"time"

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

	// parse flags
	report := flag.String("report", "", "generate a report")
	flag.Parse()

	// Initialize Dependencies
	heimatAPI := api.NewAPI("https://heimat.sprinteins.com/api/v1")

	if report != nil && *report == "times" {
		runTimeReport(heimatAPI)
		return
	}
	startPrompt(heimatAPI)

}

func runTimeReport(api *api.API) {
	emps := api.FetchEmployeeIDs()
	timeReports := make([]print.TimeReport, len(emps))
	var wg sync.WaitGroup
	wg.Add(len(emps))
	for ei, emp := range emps {
		go func(index int, emp heimat.User) {
			balances := api.FetchBalancesByUser(emp.ID, time.Now())
			timeReports[index] = print.TimeReport{Name: emp.Name(), TimeBalance: balances.BalanceWorkingHours, VacationLeft: balances.HolidayEntitlement - balances.Holidays}
			wg.Done()
		}(ei, emp)
	}
	wg.Wait()

	print.TimeReports(timeReports)
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

	// Skip Login if already authenticated
	if heimatAPI.IsAuthenticated() {
		sm.ChangeState(stateKeyHome)
	} else {
		sm.ChangeState(stateKeyLogin)
	}

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
