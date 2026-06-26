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
			label: "Should pass",
			csvString: `testLabel,label,value,expectedExecutionLabel,objectId,objectType,actionUID,actionURL,portalId
			1,Brian Halligan,1234,Success,123,CONTACT,test_action,http://localhost:3000,111
			2,Maria Johnson,9876,Failure,987,CONTACT,new_action,https://api.hubapi.com/crm,222`,
			wantErr: false,
			result: []TestCase{
				{
					TestLabel: "1",
					ActionUID: "test_action",
					ActionURL: "http://localhost:3000",
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
					TestLabel: "2",
					ActionUID: "new_action",
					ActionURL: "https://api.hubapi.com/crm",
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
			label: "Error - mismatched row length",
			csvString: `testLabel,label,value,expectedExecutionLabel
			1,Brian Halligan,1234,Success,Should Error
			2,Maria Johnson,9876,Failure`,
			wantErr: true,
		},
		{
			label:     "Error - no records",
			csvString: "testLabel,actionUID,label,value,expectedExecutionLabel",
			wantErr:   true,
		},
		{
			label: "Error - no actionUID",
			csvString: `testLabel,label,value,expectedExecutionLabel
			1,Brian Halligan,1234,Success
			2,Maria Johnson,9876,Failure`,
			wantErr: true,
		},
		{
			label: "Error - invalid objectId",
			csvString: `testLabel,label,value,expectedExecutionLabel,objectId
			1,Brian Halligan,1234,Success,asd
			2,Maria Johnson,9876,Failure,123`,
			wantErr: true,
		},
		{
			label: "Error - invalid portalId",
			csvString: `testLabel,label,value,expectedExecutionLabel,portalId
			1,Brian Halligan,1234,Success,asd
			2,Maria Johnson,9876,Failure,123`,
			wantErr: true,
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
