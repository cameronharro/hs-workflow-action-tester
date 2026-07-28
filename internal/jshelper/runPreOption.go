package jshelper

import (
	"encoding/json"
	"fmt"

	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

type PreOptionEvent struct {
	WebhookURL string `json:"webhookUrl"`
	Origin     struct {
		PortalID int `json:"portalId"`
	} `json:"origin"`
	ObjectTypeID   string                               `json:"objectTypeId"`
	InputFieldName string                               `json:"inputFieldName"`
	InputFields    map[string]testcase.OptionInputField `json:"inputFields"`
}

func (payload PreOptionEvent) getEventType() FunctionType {
	return PreActionExecution
}

type PreOptionCallback struct {
	WebhookURL  string            `json:"webhookUrl"`
	Body        any               `json:"body"`
	HttpHeaders map[string]string `json:"httpHeaders"`
	ContentType string            `json:"contentType"`
	Accept      string            `json:"accept"`
	HttpMethod  Method            `json:"httpMethod"`
}

func (callback PreOptionCallback) getCallbackType() FunctionType {
	return PreActionExecution
}

func validatePreOption(data []byte) (PreOptionCallback, error) {
	result := PreOptionCallback{}
	if err := json.Unmarshal(data, &result); err != nil {
		return PreOptionCallback{}, err
	}

	hasValidMethod := isValidMethod(result.HttpMethod)

	if !hasValidMethod {
		return PreOptionCallback{}, fmt.Errorf("[PreActionCallback]: Invalid HTTP Method %v", result.HttpMethod)
	}

	return result, nil
}

func RunPreOptionFunction(event PreOptionEvent, function string) (PreOptionCallback, error) {
	return spawn(event, function, validatePreOption)
}
