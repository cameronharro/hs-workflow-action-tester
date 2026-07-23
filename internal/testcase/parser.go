package testcase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type TestType string

const (
	Action TestType = "action"
	Option TestType = "option"
)

type TestCase struct {
	TestLabel              string         `json:"testLabel"`
	ActionUID              string         `json:"actionUID"`
	ActionURL              string         `json:"actionURL"`
	InputFields            map[string]any `json:"inputFields"`
	ObjectID               int            `json:"objectID"`
	ObjectType             string         `json:"objectType"`
	PortalID               int            `json:"portalID"`
	ExpectedExecutionLabel string         `json:"expectedExecutionLabel"`
}

func Parse(filePath string) ([]TestCase, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file at %s: %w", filePath, err)
	}
	defer file.Close()

	var testCases []TestCase
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	err = decoder.Decode(&testCases)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s: %w", filePath, err)
	}

	if len(testCases) == 0 {
		return nil, fmt.Errorf("No test cases provided at %s: %w", filePath, err)
	}

	var invalidTestCases error
	for i, testCase := range testCases {
		if testCase.ActionUID == "" {
			invalidTestCases = errors.Join(fmt.Errorf("[TestCase %d]: No actionUID provided", i))
		}
		if testCase.ObjectType == "" {
			invalidTestCases = errors.Join(fmt.Errorf("[TestCase %d]: No objectType provided", i))
		}
		if testCase.TestLabel == "" {
			invalidTestCases = errors.Join(fmt.Errorf("[TestCase %d]: No testLabel provided", i))
		}
		if testCase.PortalID == 0 {
			invalidTestCases = errors.Join(fmt.Errorf("[TestCase %d]: No portalID provided", i))
		}
		if testCase.ObjectID == 0 {
			invalidTestCases = errors.Join(fmt.Errorf("[TestCase %d]: No objectID provided", i))
		}
	}

	if invalidTestCases != nil {
		return nil, invalidTestCases
	}

	return testCases, nil
}
