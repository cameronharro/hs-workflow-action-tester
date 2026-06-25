package hsserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

type TestCaseError struct {
	testCase testcase.TestCase
	error    error
}

func (e *TestCaseError) Error() string {
	return fmt.Sprintf("[TestCase %s]: %s", e.testCase.TestLabel, e.error.Error())
}

func (s *HSServer) RunTestCase(
	testCase testcase.TestCase,
	actionDefs []actiondefinition.ActionDefinition,
) error {
	s.testsInitiated.Add(1)
	actionDef, err := getDefForCase(testCase, actionDefs)
	if err != nil {
		return &TestCaseError{testCase, err}
	}

	if err = validateCaseAgainstDef(testCase, actionDef); err != nil {
		return &TestCaseError{testCase, err}
	}

	req, callbackId, err := s.createRequest(actionDef, testCase)
	if err != nil {
		return &TestCaseError{testCase, err}
	}

	s.resolutionQueue.PayloadChan <- testPayload{
		CallbackId: callbackId,
		TestCase:   testCase,
		ActionDef:  actionDef,
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &TestCaseError{testCase, err}
	}

	if res.StatusCode >= 300 {
		return &TestCaseError{
			testCase: testCase,
			error:    fmt.Errorf("Response status %d from application", res.StatusCode),
		}
	}

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return &TestCaseError{testCase, err}
	}

	var responseJSON map[string]any
	err = json.Unmarshal(responseBytes, &responseJSON)
	if err != nil {
		return &TestCaseError{testCase, err}
	}

	s.resolutionQueue.ResponseChan <- testResponse{
		CallbackId:   callbackId,
		ResponseBody: responseJSON,
	}
	return nil
}
