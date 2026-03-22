package jshelper_test

import (
	"errors"
	"fmt"
	"os/exec"
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
			ExpectVal: jshelper.PreActionCallback{},
		},
		{
			Name:      "Infinite Loop",
			Event:     jshelper.PreActionEvent{},
			Function:  `while(true){}`,
			ExpectErr: true,
			ExpectVal: jshelper.PreActionCallback{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result, err := jshelper.RunPreActionFunction(testCase.Event, testCase.Function)
			if err != nil != testCase.ExpectErr {
				if err == nil {
					t.Fatal("Expected Error but got nil")
				}
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					fmt.Println(string(exitErr.Stderr))
				}
				t.Fatal(err.Error())
			}
			if !reflect.DeepEqual(result, testCase.ExpectVal) {
				t.Fatalf("Expected %s, got %s", testCase.ExpectVal, result)
			}
		})
	}
}
