package jshelper

import (
	"encoding/json"
)

type PostActionEvent map[string]any

func (payload PostActionEvent) getEventType() FunctionType {
	return PostActionExecution
}

type PostActionCallback struct {
	OutputFields map[string]string `json:"outputFields"`
}

func (callback PostActionCallback) getCallbackType() FunctionType {
	return PostActionExecution
}

func validatePostAction(data []byte) (PostActionCallback, error) {
	result := PostActionCallback{}
	if err := json.Unmarshal(data, &result); err != nil {
		return PostActionCallback{}, err
	}

	return result, nil
}

func RunPostActionFunction(event PostActionEvent, function string) (PostActionCallback, error) {
	return spawn(event, function, validatePostAction)
}
