package fsm

import (
	"log"
)

type State string

type StateMachine struct {
	ExistingStates []*State
	UserStates     map[string]*State
}

func NewStateMachine() *StateMachine {
	log.Println("New FSM created")
	return &StateMachine{
		UserStates: map[string]*State{},
	}
}

func (sm *StateMachine) NewState(name State) *State {
	sm.ExistingStates = append(sm.ExistingStates, &name)
	return &name
}

func (sm *StateMachine) SetState(user string, state State) {
	sm.UserStates[user] = &state
	log.Printf("%v set for %v", state, user)
}

func (sm *StateMachine) GetState(user string) *State {
	state := sm.UserStates[user]
	if state == nil {
		sm.SetState(user, InitialState)
		state = sm.UserStates[user]
	}
	return state
}
