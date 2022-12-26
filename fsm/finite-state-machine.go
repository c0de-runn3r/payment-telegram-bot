package fsm

import (
	"log"
)

type State string

type StateMachine struct {
	ExistingStates []*State
	UserStates     map[int64]*State
}

func NewStateMachine() *StateMachine {
	log.Println("New FSM created")
	return &StateMachine{
		UserStates: map[int64]*State{},
	}
}

func (sm *StateMachine) NewState(name State) *State {
	sm.ExistingStates = append(sm.ExistingStates, &name)
	return &name
}

func (sm *StateMachine) SetState(userID int64, state State) {
	sm.UserStates[userID] = &state
	log.Printf("%v set for %v", state, userID)
}

func (sm *StateMachine) GetState(userID int64) *State {
	state := sm.UserStates[userID]
	if state == nil {
		sm.SetState(userID, InitialState)
		state = sm.UserStates[userID]
	}
	return state
}
