package main

import (
	"heimatcli/x/log"

	"github.com/c-bata/go-prompt"
)

type State interface {
	Suggestions(in prompt.Document) []prompt.Suggest
	Prefix() string
	Exe(in string)
}

// StateMachine _
type StateMachine struct {
	currentState State
	states       map[string]State
}

// NewStateMachine _
func NewStateMachine() *StateMachine {
	return &StateMachine{
		states: make(map[string]State),
	}
}

// CurrentState _
func (sm StateMachine) CurrentState() State {
	return sm.currentState
}

// AddState _
func (sm *StateMachine) AddState(key string, s State) {
	sm.states[key] = s
}

// ChangeState _
func (sm *StateMachine) ChangeState(key string) {
	s, ok := sm.states[key]
	if !ok {
		log.Error.Printf("could not find state with key:%s", key)
	}
	sm.currentState = s
}
