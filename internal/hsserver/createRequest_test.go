package hsserver_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/hsserver"
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

	testCases := []TestCase{
		{
			label: "working baseline",
			actionDef: actiondefinition.ActionDefinition{
				Config: actiondefinition.ActionConfig{
					ActionURL: "http://localhost:8000/contacts",
				},
			},
			testCase: testcase.TestCase{
				InputFields: map[string]any{
					"foo": "bar",
				},
				ObjectID:   1234,
				ObjectType: "CONTACT",
				PortalID:   11,
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
			label: "mutates with function",
			actionDef: actiondefinition.ActionDefinition{
				Config: actiondefinition.ActionConfig{
					ActionURL: "http://localhost:8000/contacts",
					Functions: []actiondefinition.Function{
						actiondefinition.ActionFunction{
							FunctionType: actiondefinition.PreActionExecution,
							FunctionSource: `exports.main = (event) => {
								return {webhookUrl: "https://google.com", body: event.inputFields, httpMethod: "PATCH"}
							}`,
						},
					},
				},
			},
			testCase: testcase.TestCase{
				InputFields: map[string]any{
					"foo": "bar",
				},
				ObjectID:   1234,
				ObjectType: "CONTACT",
				PortalID:   11,
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

	server := hsserver.NewHSServer("asdfaegwagfasgrasef")
	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			req, _, err := server.CreateRequest(testCase.actionDef, testCase.testCase)
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
