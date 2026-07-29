package jshelper_test

import (
	"reflect"
	"testing"

	"github.com/cameronharro/hs-workflow-tester/internal/jshelper"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func TestRunPostOptionFunction(t *testing.T) {
	type TestCase struct {
		Name      string
		Event     jshelper.PostOptionEvent
		Function  string
		ExpectErr bool
		ExpectVal jshelper.PostOptionCallback
	}
	testCases := []TestCase{
		{
			Name: "Returns invalid options - string",
			Function: `exports.main = (_, callback) => {
				return callback({options: "string"})
			}`,
			ExpectErr: true,
		},
		{
			Name: "Returns invalid options - number",
			Function: `exports.main = (_, callback) => {
				return callback({options: 24})
			}`,
			ExpectErr: true,
		},
		{
			Name: "Returns invalid options - empty object",
			Function: `exports.main = (_, callback) => {
				return callback({options: {}})
			}`,
			ExpectErr: true,
		},
		{
			Name: "Returns no options",
			Function: `exports.main = (_, callback) => {
				return callback({})
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PostOptionCallback{},
		},
		{
			Name: "Returns valid options - empty array",
			Function: `exports.main = (_, callback) => {
				return callback({options: []})
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PostOptionCallback{Options: []testcase.Option{}},
		},
		{
			Name: "Returns valid options",
			Function: `exports.main = (_, callback) => {
				return callback({options: [{label:"Big Widget", "description": "", value: "big_widget"}]})
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PostOptionCallback{Options: []testcase.Option{{Label: "Big Widget", Value: "big_widget"}}},
		},
		{
			Name: "Transforms receieved payload",
			Event: jshelper.PostOptionEvent{
				FieldKey: "test",
				ResponseBody: map[string]any{
					"opts": []string{"block", "async", "success"},
				},
			},
			Function: `exports.main = (event, callback) => {
				const { opts } = event.responseBody
				return callback({options: opts.map(ele => ({label: ele.toUpperCase(), value: ele}))})
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PostOptionCallback{Options: []testcase.Option{
				{Label: "BLOCK", Value: "block"},
				{Label: "ASYNC", Value: "async"},
				{Label: "SUCCESS", Value: "success"},
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result, err := jshelper.RunPostOptionFunction(testCase.Event, testCase.Function)
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
