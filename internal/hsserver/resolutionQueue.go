package hsserver

import (
	"context"
	"encoding/json"
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
	allTestsBegun      *atomic.Bool
	testsInitiated     *atomic.Int64
	responsesProcessed *atomic.Int64
	payloadChan        chan<- testPayload
	responseChan       chan<- testResponse
	close              func()
}

func newResolutionQueue(server *HSServer) *resolutionQueue {
	payloadChan := make(chan testPayload)
	responseChan := make(chan testResponse)
	testsInitiated := &atomic.Int64{}
	responsesProcessed := &atomic.Int64{}
	allTestsBegun := &atomic.Bool{}
	ctx, cancelFunc := context.WithCancel(server.ctx)
	doneChan := make(chan struct{})
	closeFunc := func() {
		cancelFunc()
		<-doneChan
	}
	resolutionQueue := resolutionQueue{
		payloadChan:        payloadChan,
		responseChan:       responseChan,
		allTestsBegun:      allTestsBegun,
		testsInitiated:     testsInitiated,
		responsesProcessed: responsesProcessed,
		close:              closeFunc,
	}
	go func() {
		defer cancelFunc()
		payloads := map[string]testPayload{}

	ProcessingLoop:
		for {
			select {
			case payload := <-payloadChan:
				payloads[payload.CallbackId] = payload

			case response := <-responseChan:
				payload, ok := payloads[response.CallbackId]
				if !ok {
					server.AddResult(fmt.Errorf("Missing testCase for callback: %s", response.CallbackId))
					continue
				}

				comparisonError := checkResponseAgainstPayload(payload, response)
				if errors.Is(comparisonError, asyncSignal) {
					continue
				}
				delete(payloads, response.CallbackId)
				server.AddResult(comparisonError)
			case <-ctx.Done():
				break ProcessingLoop
			case <-time.After(1 * time.Second):
				if allTestsBegun.Load() && (ctx.Err() != nil || responsesProcessed.Load() == testsInitiated.Load()) {
					break ProcessingLoop
				}
			}
		}
		for callbackId, payload := range payloads {
			server.AddResult(
				&TestCaseError{
					testCase: payload.TestCase,
					error: fmt.Errorf(
						"No response received for callback Id %s",
						callbackId,
					),
				},
			)
		}
		if int(responsesProcessed.Load()) < int(testsInitiated.Load()) {
			server.AddResult(
				fmt.Errorf(
					"[TestCases]: Expected %d cases, received %d",
					testsInitiated.Load(),
					int(responsesProcessed.Load()),
				),
			)
		}
		server.close()
		doneChan <- struct{}{}
	}()

	return &resolutionQueue
}

var asyncSignal = errors.New("Resolution will be asynchronous, ignore")

func checkResponseAgainstPayload(payload testPayload, response testResponse) error {
	switch test := payload.TestCase.Test.(type) {
	case testcase.ActionTest:
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
			return asyncSignal
		}

		executionRuleLabel := payload.ActionDef.GetMatchingExecutionRule(outputFields)
		if executionRuleLabel != test.ExpectedExecutionRule {
			return &TestCaseError{
				testCase: payload.TestCase,
				error: fmt.Errorf(
					"Expected execution label %s, received %s",
					test.ExpectedExecutionRule,
					executionRuleLabel,
				),
			}
		}
		return nil
	case testcase.OptionTest:
		options := []testcase.Option{}
		if postOptionsFunc := payload.ActionDef.GetPostOptionFunction(test.InputFieldName); postOptionsFunc != nil {
			postOptionCallback, jsErr := jshelper.RunPostOptionFunction(
				jshelper.PostOptionEvent{
					FieldKey:     test.InputFieldName,
					ResponseBody: response.ResponseBody,
				},
				postOptionsFunc.SourceCode(),
			)
			if jsErr != nil {
				return &TestCaseError{
					testCase: payload.TestCase,
					error:    jsErr,
				}
			}
			options = postOptionCallback.Options
		} else if optsValue, ok := response.ResponseBody["options"]; ok {
			data, err := json.Marshal(optsValue)
			if err != nil {
				return &TestCaseError{
					testCase: payload.TestCase,
					error:    err,
				}
			}

			err = json.Unmarshal(data, &options)
			if err != nil {
				return &TestCaseError{
					testCase: payload.TestCase,
					error:    err,
				}
			}
		}

		if !slices.Equal(options, test.ExpectedOptions) {
			return &TestCaseError{
				testCase: payload.TestCase,
				error: fmt.Errorf(
					"Expected options %v, received %v",
					test.ExpectedOptions,
					options,
				),
			}
		}
		return nil
	}
	return fmt.Errorf("Unknown test type")
}
