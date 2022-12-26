package fsm

import "context"

// to initialize state machine + states

var SM = NewStateMachine()

var (
	InitialState       = *SM.NewState("InitialState")
	QuestionState      = *SM.NewState("QuestionState")
	ReplyQuestionState = *SM.NewState("ReplyQuestionState")
)

var CTX = context.Background()
