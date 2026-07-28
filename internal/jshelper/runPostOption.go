package jshelper

import (
	"encoding/json"
)

type PostOptionEvent struct {
	FieldKey     string `json:"fieldKey"`
	ResponseBody string `json:"responseBody"`
}

func (payload PostOptionEvent) getEventType() FunctionType {
	return PostFetchOptions
}

type Option struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

type PostOptionCallback struct {
	Options []Option
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
