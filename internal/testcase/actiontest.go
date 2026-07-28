package testcase

import (
	"bytes"
	"encoding/json"
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

func (t *ActionTest) UnmarshalJSON(data []byte) error {
	type Peek struct {
		ActionURL             string                     `json:"actionURL"`
		ObjectID              int                        `json:"objectID"`
		ObjectType            string                     `json:"objectType"`
		ExpectedExecutionRule string                     `json:"expectedExecutionRule"`
		InputFields           map[string]json.RawMessage `json:"inputFields"`
	}
	peek := Peek{}
	err := json.Unmarshal(data, &peek)
	if err != nil {
		return err
	}

	t.ActionURL = peek.ActionURL
	t.ObjectID = peek.ObjectID
	t.ObjectType = peek.ObjectType
	t.ExpectedExecutionRule = peek.ExpectedExecutionRule
	t.InputFields = map[string]any{}

	for k, v := range peek.InputFields {
		var val any
		decoder := json.NewDecoder(bytes.NewBuffer(v))
		decoder.UseNumber()
		err := decoder.Decode(&val)
		if err != nil {
			return err
		}
		t.InputFields[k] = val
	}

	return nil
}
