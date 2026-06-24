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
			Name: "Returns empty callback ",
			Function: `exports.main = () => {
				return {}
			}`,
			ExpectErr: true,
		},
		{
			Name: "Returns valid callback",
			Function: `exports.main = () => {
				return {httpMethod: "POST", body: {foo: "bar"}}
			}`,
			ExpectErr: false,
			ExpectVal: jshelper.PreActionCallback{HttpMethod: jshelper.Post, Body: map[string]any{"foo": "bar"}},
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
