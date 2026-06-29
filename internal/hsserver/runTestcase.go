package hsserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
) {
	actionDef, err := getDefForCase(testCase, actionDefs)
	if err != nil {
		s.result.Add(&TestCaseError{testCase, err})
		return
	}

	if err = validateCaseAgainstDef(testCase, actionDef); err != nil {
		s.result.Add(&TestCaseError{testCase, err})
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	req, callbackId, err := s.createRequest(ctx, actionDef, testCase)
	if err != nil {
		s.result.Add(&TestCaseError{testCase, err})
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.result.Add(&TestCaseError{testCase, err})
		return
	}

	if res.StatusCode >= 300 {
		s.result.Add(&TestCaseError{
			testCase: testCase,
			error:    fmt.Errorf("Response status %d from application", res.StatusCode),
		})
		return
	}

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.result.Add(&TestCaseError{testCase, err})
		return
	}

	var responseJSON map[string]any
	err = json.Unmarshal(responseBytes, &responseJSON)
	if err != nil {
		s.result.Add(&TestCaseError{testCase, err})
		return
	}

	s.resolutionQueue.payloadChan <- testPayload{
		CallbackId: callbackId,
		TestCase:   testCase,
		ActionDef:  actionDef,
	}

	s.resolutionQueue.responseChan <- testResponse{
		CallbackId:   callbackId,
		ResponseBody: responseJSON,
	}
}
