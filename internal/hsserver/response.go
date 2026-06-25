package hsserver

import (
	"fmt"
	"slices"
)

type response struct {
	OutputFields map[string]any `json:"outputFields"`
}

type hsExecutionState string

const (
	Async   hsExecutionState = "ASYNC"
	Block   hsExecutionState = "BLOCK"
	Failure hsExecutionState = "FAIL_CONTINUE"
	Success hsExecutionState = "SUCCESS"
)

func getExecutionState(state any) (hsExecutionState, error) {
	str, ok := state.(string)
	if !ok {
		return "", fmt.Errorf("Invalid Execution State: %v", state)
	}
	if !slices.Contains([]hsExecutionState{Async, Block, Failure, Success}, hsExecutionState(str)) {
		return "", fmt.Errorf("Invalid Execution State: %v", str)
	}
	return hsExecutionState(str), nil
}
