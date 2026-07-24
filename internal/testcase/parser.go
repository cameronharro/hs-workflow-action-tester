package testcase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

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
			invalidTestCases = errors.Join(invalidTestCases, fmt.Errorf("[TestCase %d]: No actionUID provided", i))
		}
		if testCase.TestLabel == "" {
			invalidTestCases = errors.Join(invalidTestCases, fmt.Errorf("[TestCase %d]: No testLabel provided", i))
		}
		if testCase.PortalID == 0 {
			invalidTestCases = errors.Join(invalidTestCases, fmt.Errorf("[TestCase %d]: No portalID provided", i))
		}
		if err := testCase.Test.Validate(); err != nil {
			invalidTestCases = errors.Join(invalidTestCases, fmt.Errorf("[TestCase %d]: %w", i, err))
		}
	}

	if invalidTestCases != nil {
		return nil, invalidTestCases
	}

	return testCases, nil
}
