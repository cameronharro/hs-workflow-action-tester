package jshelper_test

import (
	"reflect"
	"testing"

	"github.com/cameronharro/hs-workflow-tester/internal/jshelper"
)

func TestRunPreActionFunction(t *testing.T) {
	type TestCase struct {
		Name      string
		Event     jshelper.PreActionEvent
		Function  string
		ExpectErr bool
		ExpectVal jshelper.PreActionCallback
	}
	testCases := []TestCase{
		{
			Name:      "Hello world",
			Event:     jshelper.PreActionEvent{},
			Function:  `console.log("Hello world!")`,
			ExpectErr: true,
		},
		{
			Name:  "Infinite Loop",
			Event: jshelper.PreActionEvent{},
			Function: `exports.main = () => {
				while(true){}
			}`,
			ExpectErr: true,
		},
		{
			Name:  "Returns string",
			Event: jshelper.PreActionEvent{},
			Function: `exports.main = () => {
				return "Hello World!"
			}`,
			ExpectErr: true,
		},
		{
			Name:  "Returns null",
			Event: jshelper.PreActionEvent{},
			Function: `exports.main = () => {
				return null
			}`,
			ExpectErr: true,
		},
		{
			Name:  "Returns function",
			Event: jshelper.PreActionEvent{},
			Function: `exports.main = () => {
				function result() {
					return true
				}
				return result
			}`,
			ExpectErr: true,
		},
		{
			Name:  "Returns empty callback without erroring",
			Event: jshelper.PreActionEvent{},
			Function: `exports.main = () => {
				return {httpMethod: "POST"}
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PreActionCallback{HttpMethod: jshelper.Post},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result, err := jshelper.RunPreActionFunction(testCase.Event, testCase.Function)
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
