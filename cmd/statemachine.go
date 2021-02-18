package main

import (
	"heimatcli/x/log"

	"github.com/c-bata/go-prompt"
)

// State _
type State interface {
	Suggestions(in prompt.Document) []prompt.Suggest
	Prefix() string
	Exe(in string) StateKey
	Init()
}

// StateKey _
type StateKey string

const (
	stateKeyNoChange = ""

	stateKeyLogin   StateKey = "login"
	stateKeyHome    StateKey = "home"
	stateKeyTimeAdd StateKey = "timeadd"
)

// StateMachine _
type StateMachine struct {
	currentState State
	states       map[StateKey]State
	homeStateKey StateKey
}

// NewStateMachine _
func NewStateMachine(homeStateKey StateKey) *StateMachine {
	return &StateMachine{
		states:       make(map[StateKey]State),
		homeStateKey: homeStateKey,
	}
}

// CurrentState _
func (sm StateMachine) CurrentState() State {
	return sm.currentState
}

// AddState _
func (sm *StateMachine) AddState(key StateKey, s State) {
	sm.states[key] = s
}

// ChangeState _
func (sm *StateMachine) ChangeState(key StateKey) {
	s, ok := sm.states[key]
	if !ok {
		log.Error.Printf("could not find state with key:%s", key)
	}
	sm.currentState = s
	sm.currentState.Init()
}

// Cancel cancels the current state and returns to the home state
func (sm *StateMachine) Cancel() {
	sm.ChangeState(sm.homeStateKey)
}
