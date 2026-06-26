package hsserver

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync/atomic"
	"time"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/jshelper"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

type testPayload struct {
	CallbackId string
	TestCase   testcase.TestCase
	ActionDef  actiondefinition.ActionDefinition
}
type testResponse struct {
	CallbackId   string
	ResponseBody map[string]any
}

type resolutionQueue struct {
	PayloadChan  chan<- testPayload
	ResponseChan chan<- testResponse
}

func startResolutionQueue(testCaseCount *atomic.Int64, resultChan chan<- error, timeout time.Duration) resolutionQueue {
	payloadChan := make(chan testPayload)
	responseChan := make(chan testResponse)
	resolutionQueue := resolutionQueue{
		PayloadChan:  payloadChan,
		ResponseChan: responseChan,
	}

	payloads := map[string]testPayload{}
	responsesProcessed := 0
	var err error

	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
		defer cancelFunc()
	ProcessingLoop:
		for {
			select {
			case payload := <-payloadChan:
				payloads[payload.CallbackId] = payload

			case response := <-responseChan:

				payload, ok := payloads[response.CallbackId]
				if !ok {
					err = errors.Join(err, fmt.Errorf("Missing testCase for callback: %s", response.CallbackId))
					continue
				}

				comparisonError := checkResponseAgainstPayload(payload, response)
				if comparisonError != nil {
					var target *TestCaseError
					if !errors.As(comparisonError, &target) {
						continue
					}
					err = errors.Join(err, comparisonError)
				}

				responsesProcessed++
				delete(payloads, response.CallbackId)
				if responsesProcessed >= int(testCaseCount.Load()) {
					break ProcessingLoop
				}

			case <-ctx.Done():
				for callbackId, payload := range payloads {
					err = errors.Join(
						err,
						&TestCaseError{
							testCase: payload.TestCase,
							error: fmt.Errorf(
								"No response received for callback Id %s",
								callbackId,
							),
						},
					)
				}
				if len(payloads)+responsesProcessed < int(testCaseCount.Load()) {
					err = errors.Join(
						err,
						fmt.Errorf(
							"[TestCases]: Expected %d cases, received %d",
							testCaseCount.Load(),
							len(payloads)+responsesProcessed,
						),
					)
				}
				break ProcessingLoop
			}
		}
		resultChan <- err
	}()

	return resolutionQueue
}

func checkResponseAgainstPayload(payload testPayload, response testResponse) error {
	var outputFields map[string]any
	if postActionFunc := payload.ActionDef.GetPostActionFunction(); postActionFunc != nil {
		postActionCallback, jsErr := jshelper.RunPostActionFunction(
			response.ResponseBody,
			postActionFunc.SourceCode(),
		)
		if jsErr != nil {
			return &TestCaseError{
				testCase: payload.TestCase,
				error:    jsErr,
			}
		}
		outputFields = postActionCallback.OutputFields
	} else if rawOutputFields, ok := response.ResponseBody["outputFields"]; ok {
		outputFields, ok = rawOutputFields.(map[string]any)
	}

	hsExState, stateErr := getExecutionState(outputFields["hs_execution_state"])
	if stateErr != nil {
		return &TestCaseError{
			testCase: payload.TestCase,
			error:    stateErr,
		}
	}
	if slices.Contains([]hsExecutionState{Async, Block}, hsExState) {
		return fmt.Errorf("Resolution will be asynchronous, ignore")
	}

	executionRuleLabel := payload.ActionDef.GetMatchingExecutionRule(outputFields)
	if executionRuleLabel != payload.TestCase.ExpectedExecutionLabel {
		return &TestCaseError{
			testCase: payload.TestCase,
			error: fmt.Errorf(
				"Expected execution label %s, received %s",
				payload.TestCase.ExpectedExecutionLabel,
				executionRuleLabel,
			),
		}
	}
	return nil
}
