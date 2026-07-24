package testcase

import (
	"errors"
)

type ActionTest struct {
	ActionURL             string         `json:"actionURL"`
	InputFields           map[string]any `json:"inputFields"`
	ObjectID              int            `json:"objectID"`
	ObjectType            string         `json:"objectType"`
	ExpectedExecutionRule string         `json:"expectedExecutionRule"`
}

func (t ActionTest) Validate() error {
	var result error
	if t.ObjectType == "" {
		result = errors.Join(result, errors.New("Missing objectType"))
	}

	if t.ObjectID == 0 {
		result = errors.Join(result, errors.New("Missing objectID"))
	}

	return result
}
