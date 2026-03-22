package jshelper

import (
	"encoding/json"
	"fmt"
	"slices"
)

type PreActionEvent struct {
}

func (payload PreActionEvent) getEventType() FunctionType {
	return PreActionExecution
}

type PreActionCallback struct {
	WebhookURL  string            `json:"webhookUrl"`
	Body        string            `json:"body"`
	HttpHeaders map[string]string `json:"httpHeaders"`
	ContentType string            `json:"contentType"`
	Accept      string            `json:"accept"`
	HttpMethod  Method            `json:"httpMethod"`
}

func (callback PreActionCallback) getCallbackType() FunctionType {
	return PreActionExecution
}

func isValidMethod(str Method) bool {
	return slices.Contains(getAllMethods(), str)
}

func validatePreAction(data []byte) (PreActionCallback, error) {
	result := PreActionCallback{}
	if err := json.Unmarshal(data, &result); err != nil {
		return PreActionCallback{}, err
	}

	hasValidMethod := isValidMethod(result.HttpMethod)

	if !hasValidMethod {
		return PreActionCallback{}, fmt.Errorf("Invalid PreActionCallback: %v", result)
	}

	return result, nil
}

func RunPreActionFunction(event PreActionEvent, function string) (PreActionCallback, error) {
	return spawn(event, function, validatePreAction)
}
