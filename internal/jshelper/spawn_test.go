package jshelper_test

import (
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"testing"

	"github.com/cameronharro/hs-workflow-tester/internal/jshelper"
)

func TestSpawn(t *testing.T) {
	type TestCase struct {
		Name      string
		Data      jshelper.RequestParams
		Code      string
		ExpectErr bool
		ExpectVal jshelper.RequestParams
	}
	testCases := []TestCase{
		{
			Name:      "Hello world",
			Data:      jshelper.RequestParams{},
			Code:      `console.log("Hello world!")`,
			ExpectErr: false,
			ExpectVal: jshelper.RequestParams{},
		},
		{
			Name:      "Infinite Loop",
			Data:      jshelper.RequestParams{},
			Code:      `while(true){}`,
			ExpectErr: true,
			ExpectVal: jshelper.RequestParams{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result, err := jshelper.Spawn(testCase.Data, testCase.Code)
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
