package hsserver

import (
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func TestResolutionQueue(t *testing.T) {
	type ReqResCycle struct {
		payload   testPayload
		responses []testResponse
	}
	type TestCase struct {
		label            string
		timeout          time.Duration
		startingReqCount int
		reqResCycles     []ReqResCycle
		expectErr        bool
	}

	actionDef := actiondefinition.ActionDefinition{
		Config: actiondefinition.ActionConfig{
			ExecutionRules: []actiondefinition.ExecutionRule{
				actiondefinition.ExecutionRule{
					LabelName: "Success",
					Conditions: map[string]any{
						"status": "success",
					},
				},
				actiondefinition.ExecutionRule{
					LabelName: "Failure",
					Conditions: map[string]any{
						"status": "failure",
					},
				},
			},
		},
	}
	actionDefWithPostActionFunction := actionDef
	actionDefWithPostActionFunction.Config.Functions = []actiondefinition.Function{
		actiondefinition.ActionFunction{
			FunctionType:   actiondefinition.PostActionExecution,
			FunctionSource: `exports.main = function(event, callback) { return callback({outputFields:{hs_execution_state:"SUCCESS",status:"success"}})}`,
		},
	}
	actionDefWithInvalidPostActionFunction := actionDef
	actionDefWithInvalidPostActionFunction.Config.Functions = []actiondefinition.Function{
		actiondefinition.ActionFunction{
			FunctionType:   actiondefinition.PostActionExecution,
			FunctionSource: `exports.main = function(event, callback) { return {outputFields:{status:"success"}}}`,
		},
	}

	actionTestCases := []TestCase{
		TestCase{
			label:   "baseline working",
			timeout: 5 * time.Second,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.ActionTest{
								ExpectedExecutionRule: "Success",
							},
						},
						ActionDef: actionDef,
					},
					responses: []testResponse{
						{
							CallbackId: "123",
							ResponseBody: map[string]any{
								"outputFields": map[string]any{
									"hs_execution_state": "SUCCESS",
									"status":             "success",
								},
							},
						},
					},
				},
				{
					payload: testPayload{
						CallbackId: "456",
						TestCase: testcase.TestCase{
							TestLabel: "failure",
							Test: testcase.ActionTest{
								ExpectedExecutionRule: "Failure",
							},
						},
						ActionDef: actionDef,
					},
					responses: []testResponse{
						testResponse{
							CallbackId: "456",
							ResponseBody: map[string]any{
								"outputFields": map[string]any{
									"hs_execution_state": "FAIL_CONTINUE",
									"status":             "failure",
								},
							},
						},
					},
				},
				{
					payload: testPayload{
						CallbackId: "789",
						TestCase: testcase.TestCase{
							TestLabel: "blank",
							Test: testcase.ActionTest{
								ExpectedExecutionRule: "",
							},
						},
						ActionDef: actionDef,
					},
					responses: []testResponse{
						testResponse{
							CallbackId: "789",
							ResponseBody: map[string]any{
								"outputFields": map[string]any{
									"hs_execution_state": "SUCCESS",
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		TestCase{
			label:   "error - no response",
			timeout: 250 * time.Millisecond,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.ActionTest{
								ExpectedExecutionRule: "Success",
							},
						},
						ActionDef: actionDef,
					},
				},
			},
			expectErr: true,
		},
		TestCase{
			label:   "error - extra response",
			timeout: 5 * time.Second,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.ActionTest{
								ExpectedExecutionRule: "Success",
							},
						},
						ActionDef: actionDef,
					},
					responses: []testResponse{
						testResponse{
							CallbackId: "456",
							ResponseBody: map[string]any{
								"outputFields": map[string]any{
									"status": "failure",
								},
							},
						},
						{
							CallbackId: "123",
							ResponseBody: map[string]any{
								"outputFields": map[string]any{
									"status": "success",
								},
							},
						},
					},
				},
			},
			expectErr: true,
		},
		TestCase{
			label:   "successful blocking response",
			timeout: 5 * time.Second,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.ActionTest{
								ExpectedExecutionRule: "Success",
							},
						},
						ActionDef: actionDef,
					},
					responses: []testResponse{
						testResponse{
							CallbackId: "123",
							ResponseBody: map[string]any{
								"outputFields": map[string]any{
									"hs_execution_state": "BLOCK",
								},
							},
						},
						{
							CallbackId: "123",
							ResponseBody: map[string]any{
								"outputFields": map[string]any{
									"hs_execution_state": "SUCCESS",
									"status":             "success",
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		TestCase{
			label:   "error - incomplete blocking response",
			timeout: 250 * time.Millisecond,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.ActionTest{
								ExpectedExecutionRule: "Success",
							},
						},
						ActionDef: actionDef,
					},
					responses: []testResponse{
						testResponse{
							CallbackId: "123",
							ResponseBody: map[string]any{
								"outputFields": map[string]any{
									"hs_execution_state": "BLOCK",
								},
							},
						},
					},
				},
			},
			expectErr: true,
		},
		TestCase{
			label:   "basic POST_ACTION_EXECUTION function",
			timeout: 250 * time.Millisecond,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.ActionTest{
								ExpectedExecutionRule: "Success",
							},
						},
						ActionDef: actionDefWithPostActionFunction,
					},
					responses: []testResponse{
						testResponse{
							CallbackId: "123",
						},
					},
				},
			},
			expectErr: false,
		},
		TestCase{
			label:   "error - invalid POST_ACTION_EXECUTION function",
			timeout: 250 * time.Millisecond,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.ActionTest{
								ExpectedExecutionRule: "Success",
							},
						},
						ActionDef: actionDefWithInvalidPostActionFunction,
					},
					responses: []testResponse{
						testResponse{
							CallbackId: "123",
						},
					},
				},
			},
			expectErr: true,
		},
	}

	actionDefWithPostOptionFunction := actionDef
	actionDefWithPostOptionFunction.Config.Functions = []actiondefinition.Function{
		actiondefinition.OptionFunction{
			FunctionType:   actiondefinition.PostFetchOptions,
			Id:             "options_input",
			FunctionSource: `exports.main = function(event, callback) { return callback({options:[{label: "Test", value: "test"}]})}`,
		},
	}
	actionDefWithInvalidPostOptionFunction := actionDef
	actionDefWithInvalidPostOptionFunction.Config.Functions = []actiondefinition.Function{
		actiondefinition.OptionFunction{
			FunctionType:   actiondefinition.PostFetchOptions,
			Id:             "options_input",
			FunctionSource: `exports.main = function(event, callback) { return {options:[{label: "Test", value: "test"}]}}`,
		},
	}

	optionTestCases := []TestCase{
		TestCase{
			label:   "options - baseline working",
			timeout: 5 * time.Second,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.OptionTest{
								ExpectedOptions: []testcase.Option{
									{Label: "Test", Value: "test"},
								},
							},
						},
						ActionDef: actionDef,
					},
					responses: []testResponse{
						{
							CallbackId: "123",
							ResponseBody: map[string]any{
								"options": []map[string]any{
									map[string]any{
										"label": "Test",
										"value": "test",
									},
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		TestCase{
			label:   "options - error - bad response",
			timeout: 5 * time.Second,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.OptionTest{
								ExpectedOptions: []testcase.Option{
									{Label: "Test", Value: "test"},
								},
							},
						},
						ActionDef: actionDef,
					},
					responses: []testResponse{
						{
							CallbackId: "123",
							ResponseBody: map[string]any{
								"options": []map[string]any{
									map[string]any{
										"label": "Test Failure",
										"value": "test_failure",
									},
								},
							},
						},
					},
				},
			},
			expectErr: true,
		},
		TestCase{
			label:   "options - error - no response",
			timeout: 250 * time.Millisecond,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.OptionTest{
								ExpectedOptions: []testcase.Option{
									{Label: "Test", Value: "test"},
								},
							},
						},
						ActionDef: actionDef,
					},
				},
			},
			expectErr: true,
		},
		TestCase{
			label:   "options - basic POST_OPTION_EXECUTION function",
			timeout: 250 * time.Millisecond,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.OptionTest{
								InputFieldName: "options_input",
								ExpectedOptions: []testcase.Option{
									{Label: "Test", Value: "test"},
								},
							},
						},
						ActionDef: actionDefWithPostOptionFunction,
					},
					responses: []testResponse{
						{
							CallbackId:   "123",
							ResponseBody: map[string]any{},
						},
					},
				},
			},
			expectErr: false,
		},
		TestCase{
			label:   "options - error - invalid POST_OPTION_EXECUTION function",
			timeout: 250 * time.Millisecond,
			reqResCycles: []ReqResCycle{
				{
					payload: testPayload{
						CallbackId: "123",
						TestCase: testcase.TestCase{
							TestLabel: "success",
							Test: testcase.OptionTest{
								InputFieldName: "options_input",
								ExpectedOptions: []testcase.Option{
									{Label: "Test", Value: "test"},
								},
							},
						},
						ActionDef: actionDefWithInvalidPostOptionFunction,
					},
					responses: []testResponse{
						{
							CallbackId: "123",
							ResponseBody: map[string]any{
								"options": []map[string]any{},
							},
						},
					},
				},
			},
			expectErr: true,
		},
	}

	for _, testCase := range slices.Concat(actionTestCases, optionTestCases) {
		t.Run(testCase.label, func(t *testing.T) {
			server := NewHSServer("asdfg", rand.Intn(49151-1024)+1024, 1*time.Second)
			queue := server.resolutionQueue
			for _, cycle := range testCase.reqResCycles {
				go func() {
					queue.payloadChan <- cycle.payload
					for _, response := range cycle.responses {
						queue.responseChan <- response
					}
				}()
			}
			err := server.Wait()
			if err != nil != testCase.expectErr {
				t.Errorf("[Error]: expected? %t, got:\n%v\n", testCase.expectErr, err)
			}
		})
	}
}
