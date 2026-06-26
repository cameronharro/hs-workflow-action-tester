package hsserver

import (
	"fmt"
	"slices"
)

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
		return "", fmt.Errorf("[hs_execution_state]: Invalid state: %v", state)
	}
	if !slices.Contains([]hsExecutionState{Async, Block, Failure, Success}, hsExecutionState(str)) {
		return "", fmt.Errorf("[hs_execution_state]: Invalid state: %v", str)
	}
	return hsExecutionState(str), nil
}
