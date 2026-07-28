package testcase

import (
	"encoding/json"
	"fmt"
)

type TestCase struct {
	TestLabel string   `json:"testLabel"`
	TestType  TestType `json:"type"`
	ActionUID string   `json:"actionUID"`
	PortalID  int      `json:"portalID"`
	Test      Test     `json:"test"`
}

func (t *TestCase) UnmarshalJSON(data []byte) error {
	type Peek struct {
		TestLabel string          `json:"testLabel"`
		TestType  TestType        `json:"type"`
		ActionUID string          `json:"actionUID"`
		PortalID  int             `json:"portalID"`
		Test      json.RawMessage `json:"test"`
	}
	peek := Peek{}
	err := json.Unmarshal(data, &peek)
	if err != nil {
		return err
	}

	t.TestLabel = peek.TestLabel
	t.TestType = peek.TestType
	t.ActionUID = peek.ActionUID
	t.PortalID = peek.PortalID

	var Test Test
	switch peek.TestType {
	case Action:
		actionTest := ActionTest{}
		err = json.Unmarshal(peek.Test, &actionTest)
		if err != nil {
			return err
		}
		Test = actionTest
	case Option:
		optionTest := OptionTest{}
		err = json.Unmarshal(peek.Test, &optionTest)
		if err != nil {
			return err
		}
		Test = optionTest
	default:
		return fmt.Errorf("Unknown Test Type: %s", peek.TestType)
	}

	t.Test = Test
	return nil
}

type TestType string

const (
	Action TestType = "action"
	Option TestType = "option"
)

type Test interface {
	Validate() error
}
