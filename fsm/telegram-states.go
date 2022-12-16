package fsm

// to initialize state machine + states

var SM = NewStateMachine()

var (
	InitialState = *SM.NewState("InitialState")
)
