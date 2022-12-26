package fsm

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	sm := NewStateMachine()
	var usr int64 = 12345678

	helloState := sm.NewState("helloState")
	sm.SetState(usr, *helloState)
	fmt.Printf("%+v\n", sm.UserStates)
	fmt.Printf("user state: %+v\n", *sm.UserStates[usr])
	fmt.Printf("---->>> %+v\n", *sm.GetState(usr))
	if *sm.GetState(usr) == *helloState {
		fmt.Println("YEEEEEEESSSS")
	}

	byeState := sm.NewState("byeState")
	sm.SetState(usr, *byeState)
	fmt.Printf("%+v\n", sm.UserStates)
	fmt.Printf("user state: %+v\n", *sm.UserStates[usr])
	fmt.Printf("---->>> %+v\n", *sm.GetState(usr))
	if *sm.GetState(usr) == *byeState {
		fmt.Println("YEEEEEEESSSS X2")
	}
	fmt.Printf("User's current state is: %s\n", *sm.GetState(12345678))
}
