package jshelper

import (
	"encoding/json"

	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

type PostOptionEvent struct {
	FieldKey     string         `json:"fieldKey"`
	ResponseBody map[string]any `json:"responseBody"`
}

func (payload PostOptionEvent) getEventType() FunctionType {
	return PostFetchOptions
}

type PostOptionCallback struct {
	Options []testcase.Option
}

func (callback PostOptionCallback) getCallbackType() FunctionType {
	return PostFetchOptions
}

func validatePostOption(data []byte) (PostOptionCallback, error) {
	result := PostOptionCallback{}
	if err := json.Unmarshal(data, &result); err != nil {
		return PostOptionCallback{}, err
	}

	return result, nil
}

func RunPostOptionFunction(event PostOptionEvent, function string) (PostOptionCallback, error) {
	return spawn(event, function, validatePostOption)
}
