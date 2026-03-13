package testcases

import (
	"bytes"
	"slices"
	"testing"
)

func TestParse(t *testing.T) {
	type Case struct {
		label     string
		csvString string
		wantErr   bool
		result    []TestCase
	}
	cases := []Case{
		{
			label:     "Should pass",
			csvString: "label,value,expectedExecutionLabel\nBrian Halligan,1234,Success\nMaria Johnson,9876,Failure",
			wantErr:   false,
			result: []TestCase{
				{
					[]InputField{
						{
							name:  "label",
							value: "Brian Halligan",
						},
						{
							name:  "value",
							value: "1234",
						},
					},
					"Success",
				},
				{
					[]InputField{
						{
							name:  "label",
							value: "Maria Johnson",
						},
						{
							name:  "value",
							value: "9876",
						},
					},
					"Failure",
				},
			},
		},
		{
			label:     "Error - mismatched row length",
			csvString: "label,value,expectedExecutionLabel\nBrian Halligan,1234,Success,Should Error\nMaria Johnson,9876,Failure",
			wantErr:   true,
		},
		{
			label:     "Error - no records",
			csvString: "label,value,expectedExecutionLabel",
			wantErr:   true,
		},
	}

	for _, thisCase := range cases {
		t.Run(thisCase.label, func(t *testing.T) {
			result, err := parse(bytes.NewBuffer([]byte(thisCase.csvString)))
			if err != nil != thisCase.wantErr {
				t.Error(err.Error())
				return
			}
			matches := slices.EqualFunc(result, thisCase.result, func(a, b TestCase) bool {
				return a.ExpectedExecutionLabel == b.ExpectedExecutionLabel && slices.Equal(a.InputFields, b.InputFields)
			})
			if !matches {
				t.Errorf("Expected %v, got %v", thisCase.result, result)
			}
		})
	}
}
