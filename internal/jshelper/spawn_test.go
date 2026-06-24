package jshelper_test

import (
	"reflect"
	"testing"

	"github.com/cameronharro/hs-workflow-tester/internal/jshelper"
)

func TestSpawn(t *testing.T) {
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
			Function:  `console.log("Hello world!")`,
			ExpectErr: true,
		},
		{
			Name: "Infinite Loop",
			Function: `exports.main = () => {
				while(true){}
			}`,
			ExpectErr: true,
		},
		{
			Name: "Returns string",
			Function: `exports.main = () => {
				return "Hello World!"
			}`,
			ExpectErr: true,
		},
		{
			Name: "Returns null",
			Function: `exports.main = () => {
				return null
			}`,
			ExpectErr: true,
		},
		{
			Name: "Returns function",
			Function: `exports.main = () => {
				function result() {
					return true
				}
				return result
			}`,
			ExpectErr: true,
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
