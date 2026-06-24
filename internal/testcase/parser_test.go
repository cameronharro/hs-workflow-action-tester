package testcase

import (
	"bytes"
	"maps"
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
			csvString: "label,value,expectedExecutionLabel,objectId,actionUID,portalId\nBrian Halligan,1234,Success,123,test_action,111\nMaria Johnson,9876,Failure,987,new_action,222",
			wantErr:   false,
			result: []TestCase{
				{
					ActionUID: "test_action",
					InputFields: map[string]any{
						"label": "Brian Halligan",
						"value": "1234",
					},
					ObjectID:               123,
					ObjectType:             "CONTACT",
					PortalID:               111,
					ExpectedExecutionLabel: "Success",
				},
				{
					ActionUID: "new_action",
					InputFields: map[string]any{
						"label": "Maria Johnson",
						"value": "9876",
					},
					ObjectID:               987,
					ObjectType:             "CONTACT",
					PortalID:               222,
					ExpectedExecutionLabel: "Failure",
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
		{
			label:     "Error - no actionUID",
			csvString: "label,value,expectedExecutionLabel\nBrian Halligan,1234,Success\nMaria Johnson,9876,Failure",
			wantErr:   true,
		},
		{
			label:     "Error - invalid objectId",
			csvString: "label,value,expectedExecutionLabel,objectId\nBrian Halligan,1234,Success,asd\nMaria Johnson,9876,Failure,123",
			wantErr:   true,
		},
		{
			label:     "Error - invalid portalId",
			csvString: "label,value,expectedExecutionLabel,portalId\nBrian Halligan,1234,Success,asd\nMaria Johnson,9876,Failure,123",
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
				return a.ExpectedExecutionLabel == b.ExpectedExecutionLabel && maps.Equal(a.InputFields, b.InputFields)
			})
			if !matches {
				t.Errorf("Expected %v, got %v", thisCase.result, result)
			}
		})
	}
}
