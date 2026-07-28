package hsserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func TestCreateRequest(t *testing.T) {
	type TestCase struct {
		label       string
		actionDef   actiondefinition.ActionDefinition
		testCase    testcase.TestCase
		expectErr   bool
		expectedReq *http.Request
	}

	actionTestCases := []TestCase{
		{
			label: "actions - working baseline",
			actionDef: actiondefinition.ActionDefinition{
				Config: actiondefinition.ActionConfig{
					ActionURL: "http://localhost:8000/contacts",
				},
			},
			testCase: testcase.TestCase{
				Test: testcase.ActionTest{
					InputFields: map[string]any{
						"foo": "bar",
					},
					ObjectID:   1234,
					ObjectType: "CONTACT",
				},
				PortalID: 11,
			},
			expectErr: false,
			expectedReq: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Host:   "localhost:8000",
					Path:   "/contacts",
					Scheme: "http",
				},
				Body: io.NopCloser(bytes.NewBufferString(`{"fields":{"foo":"bar"},"object":{"objectId":1234,"objectType":"CONTACT"},"origin":{"portalId":11}}`)),
			},
		},
		{
			label: "actions - mutates with function",
			actionDef: actiondefinition.ActionDefinition{
				Config: actiondefinition.ActionConfig{
					ActionURL: "http://localhost:8000/contacts",
					Functions: []actiondefinition.Function{
						actiondefinition.ActionFunction{
							FunctionType: actiondefinition.PreActionExecution,
							FunctionSource: `exports.main = (event, callback) => {
								return callback({webhookUrl: "https://google.com", body: event.inputFields, httpMethod: "PATCH"})
							}`,
						},
					},
				},
			},
			testCase: testcase.TestCase{
				Test: testcase.ActionTest{
					InputFields: map[string]any{
						"foo": "bar",
					},
					ObjectID:   1234,
					ObjectType: "CONTACT",
				},
				PortalID: 11,
			},
			expectErr: false,
			expectedReq: &http.Request{
				Method: "PATCH",
				URL: &url.URL{
					Host:   "google.com",
					Scheme: "https",
				},
				Body: io.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
			},
		},
	}

	optionTestCases := []TestCase{
		{
			label: "options - working baseline",
			actionDef: actiondefinition.ActionDefinition{
				Config: actiondefinition.ActionConfig{
					ActionURL: "http://localhost:8000/contacts",
				},
			},
			testCase: testcase.TestCase{
				Test: testcase.OptionTest{
					InputFieldName: "custom_enum",
					InputFields: map[string]testcase.OptionInputField{
						"foo": testcase.StaticValueInputField{
							FieldType: "STATIC_VALUE",
							Value:     "bar",
						},
					},
					ObjectTypeID: "0-1",
				},
				PortalID: 11,
			},
			expectErr: false,
			expectedReq: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Host:   "localhost:8000",
					Path:   "/contacts",
					Scheme: "http",
				},
				Body: io.NopCloser(bytes.NewBufferString(`{"inputFieldName":"custom_enum","objectTypeId":"0-1","origin":{"portalId":11},"inputFields":{"foo":{"type":"STATIC_VALUE","value":"bar"}}}`)),
			},
		},
		{
			label: "options - mutates with function",
			actionDef: actiondefinition.ActionDefinition{
				Config: actiondefinition.ActionConfig{
					ActionURL: "http://localhost:8000/contacts",
					Functions: []actiondefinition.Function{
						actiondefinition.OptionFunction{
							FunctionType: actiondefinition.PreFetchOptions,
							Id:           "custom_enum",
							FunctionSource: `exports.main = (event, callback) => {
								return callback({webhookUrl: "https://google.com", body: event.inputFields, httpMethod: "PATCH"})
							}`,
						},
					},
				},
			},
			testCase: testcase.TestCase{
				Test: testcase.OptionTest{
					InputFieldName: "custom_enum",
					InputFields: map[string]testcase.OptionInputField{
						"foo": testcase.StaticValueInputField{
							FieldType: "STATIC_VALUE",
							Value:     "bar",
						},
					},
					ObjectTypeID: "0-1",
				},
				PortalID: 11,
			},
			expectErr: false,
			expectedReq: &http.Request{
				Method: "PATCH",
				URL: &url.URL{
					Host:   "google.com",
					Scheme: "https",
				},
				Body: io.NopCloser(bytes.NewBufferString(`{"foo":{"type":"STATIC_VALUE","value":"bar"}}`)),
			},
		},
	}

	server := NewHSServer("asdfaegwagfasgrasef", 8080, 1*time.Second)
	for _, testCase := range slices.Concat(actionTestCases, optionTestCases) {
		t.Run(testCase.label, func(t *testing.T) {
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
			defer cancelFunc()
			req, _, err := server.createRequest(ctx, testCase.actionDef, testCase.testCase)
			if err != nil != testCase.expectErr {
				t.Fatalf("Errors: expected? %t, got %v", testCase.expectErr, err)
			}
			if err := requestsMatch(testCase.expectedReq, req); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func requestsMatch(expected, received *http.Request) error {
	if received == nil {
		return fmt.Errorf("Missing request: expected %v, received: %v", expected, received)
	}

	var err error = nil

	type Check struct {
		failCondition bool
		message       string
	}
	eval := func(check Check, err error) error {
		if check.failCondition {
			err = errors.Join(errors.New(check.message))
		}
		return err
	}
	var expectedBody map[string]any
	expectedBodyBytes, _ := io.ReadAll(expected.Body)
	json.Unmarshal(expectedBodyBytes, &expectedBody)
	var receivedBody map[string]any
	receivedBodyBytes, _ := io.ReadAll(received.Body)
	json.Unmarshal(receivedBodyBytes, &receivedBody)
	delete(receivedBody, "callbackId")

	checks := []Check{
		{
			failCondition: expected.Method != received.Method,
			message:       fmt.Sprintf("[Method] expected %s, received %s", expected.Method, received.Method),
		},
		{
			failCondition: expected.URL.String() != received.URL.String(),
			message:       fmt.Sprintf("[URL] expected %v, received %v", expected.URL, received.URL),
		},
		{
			failCondition: !reflect.DeepEqual(expectedBody, receivedBody),
			message:       fmt.Sprintf("[Body] expected %v, received %v", expectedBody, receivedBody),
		},
	}
	for _, check := range checks {
		err = eval(check, err)
	}
	return err
}
