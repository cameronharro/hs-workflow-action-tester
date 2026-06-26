package jshelper_test

import (
	"reflect"
	"testing"

	"github.com/cameronharro/hs-workflow-tester/internal/jshelper"
)

func TestRunPostActionFunction(t *testing.T) {
	type TestCase struct {
		Name      string
		Event     jshelper.PostActionEvent
		Function  string
		ExpectErr bool
		ExpectVal jshelper.PostActionCallback
	}
	testCases := []TestCase{
		{
			Name: "Returns invalid outputFields - string",
			Function: `exports.main = (_, callback) => {
				return callback({outputFields: "string"})
			}`,
			ExpectErr: true,
		},
		{
			Name: "Returns invalid outputFields - number",
			Function: `exports.main = (_, callback) => {
				return callback({outputFields: 24})
			}`,
			ExpectErr: true,
		},
		{
			Name: "Returns no outputFields",
			Function: `exports.main = (_, callback) => {
				return callback({})
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PostActionCallback{},
		},
		{
			Name: "Returns valid outputFields - empty object",
			Function: `exports.main = (_, callback) => {
				return callback({outputFields: {}})
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PostActionCallback{OutputFields: map[string]any{}},
		},
		{
			Name: "Returns valid outputFields",
			Function: `exports.main = (_, callback) => {
				return callback({outputFields: {status: "success"}})
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PostActionCallback{OutputFields: map[string]any{"status": "success"}},
		},
		{
			Name:  "Transforms receieved payload",
			Event: jshelper.PostActionEvent{"status": 401},
			Function: `exports.main = (event, callback) => {
				return callback({outputFields: {status: event.status < 300 ? "success" : "failure"}})
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PostActionCallback{OutputFields: map[string]any{"status": "failure"}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result, err := jshelper.RunPostActionFunction(testCase.Event, testCase.Function)
			if err != nil != testCase.ExpectErr {
				if err == nil {
					t.Fatal("Expected Error but got nil")
				}
				t.Fatal(err.Error())
			}
			if !reflect.DeepEqual(result, testCase.ExpectVal) {
				t.Fatalf("Expected %s, got %s", testCase.ExpectVal, result)
			}
		})
	}
}
