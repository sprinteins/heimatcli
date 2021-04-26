package main

import (
	"flag"
	"fmt"
	"heimatcli/src/heimat"
	"heimatcli/src/heimat/api"
	"heimatcli/src/heimat/print"
	"sort"
	"strings"
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

	const defaultAPI = "https://heimat.sprinteins.com/api/v1"

	// parse flags
	report := flag.String("report", "", "generate a report")
	psl := flag.String("psl", "", "people success lead")
	apiEndpoint := flag.String("api", defaultAPI, "API Endpoint")
	versionRequest := flag.Bool("version", false, "Prints current version")

	flag.Parse()

	if *versionRequest {
		fmt.Println("v0.1.4")
		return
	}

	// Initialize Dependencies
	heimatAPI := api.NewAPI(*apiEndpoint)

	if report != nil && *report == "times" {
		cliLogin(heimatAPI)
		runTimeReport(heimatAPI, *psl)
		return
	}
	startPrompt(heimatAPI)

}

func runTimeReport(api *api.API, psl string) {
	emps := api.FetchEmployees()

	// fetch balances
	timeReports := make([]print.TimeReport, len(emps))
	var wg sync.WaitGroup
	wg.Add(len(emps))
	for ei, emp := range emps {
		go func(ei int, emp heimat.User) {
			balances := api.FetchBalancesByUser(emp.ID, time.Now())
			user := api.FetchUserByID(emp.ID)
			timeReports[ei] = print.TimeReport{
				PSL:          user.PSL,
				Name:         emp.Name(),
				TimeBalance:  balances.BalanceWorkingHours,
				VacationLeft: balances.HolidayEntitlement - balances.Holidays,
			}
			wg.Done()
		}(ei, emp)
	}
	wg.Wait()

	// print Time Reports pro PSL
	buckets := splitTimeReportToPSLBuckets(timeReports)

	fmt.Println()
	fmt.Printf("Times-Report generated at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Sorting PSLs for a consistence report,
	// because the order in ranging over map is always different
	psls := make([]string, 0, len(buckets))
	for keyPSL := range buckets {
		psls = append(psls, keyPSL)
	}
	sort.Strings(psls)

	// printing buckets
	for _, keyPSL := range psls {
		bucket := buckets[keyPSL]

		if !strings.Contains(strings.ToLower(keyPSL), strings.ToLower(psl)) {
			continue
		}
		fmt.Printf("PSL: %s\n\n", keyPSL)
		print.TimeReports(bucket)
		fmt.Printf("\n\n\n")
	}

}

type pslBuckets = map[string][]print.TimeReport

func splitTimeReportToPSLBuckets(trs []print.TimeReport) pslBuckets {
	buckets := make(pslBuckets)

	for _, tr := range trs {
		psl := tr.PSL.Name()
		_, ok := buckets[psl]
		if !ok {
			buckets[psl] = make([]print.TimeReport, 0)
		}

		bucket, _ := buckets[psl]
		buckets[psl] = append(bucket, tr)
	}

	return buckets
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
