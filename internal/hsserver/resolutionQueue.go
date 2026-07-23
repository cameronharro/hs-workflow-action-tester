package hsserver

import (
	"context"
	"errors"
	"fmt"
	"slices"
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
	payloadChan  chan<- testPayload
	responseChan chan<- testResponse
}

func newResolutionQueue(server *HSServer) *resolutionQueue {
	payloadChan := make(chan testPayload)
	responseChan := make(chan testResponse)
	resolutionQueue := resolutionQueue{
		payloadChan:  payloadChan,
		responseChan: responseChan,
	}

	go func() {
		ctx, cancelFunc := context.WithCancel(server.ctx)
		defer cancelFunc()
		payloads := map[string]testPayload{}
		testsInitiated := 0
		responsesProcessed := 0
		var err error

	ProcessingLoop:
		for {
			select {
			case payload := <-payloadChan:
				testsInitiated++
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
				if responsesProcessed >= testsInitiated {
					break ProcessingLoop
				}
			case <-time.After(1 * time.Second):
				if responsesProcessed == testsInitiated {
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
				if len(payloads)+responsesProcessed < testsInitiated {
					err = errors.Join(
						err,
						fmt.Errorf(
							"[TestCases]: Expected %d cases, received %d",
							testsInitiated,
							len(payloads)+responsesProcessed,
						),
					)
				}
				break ProcessingLoop
			}
		}
		server.result.Add(err)
		server.close()
	}()

	return &resolutionQueue
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
	if executionRuleLabel != payload.TestCase.ExpectedExecutionRule {
		return &TestCaseError{
			testCase: payload.TestCase,
			error: fmt.Errorf(
				"Expected execution label %s, received %s",
				payload.TestCase.ExpectedExecutionRule,
				executionRuleLabel,
			),
		}
	}
	return nil
}
