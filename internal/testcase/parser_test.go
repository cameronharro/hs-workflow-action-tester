package testcase_test

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func TestParse(t *testing.T) {
	type Case struct {
		label    string
		filePath string
		wantErr  bool
		result   []testcase.TestCase
	}

	cases := []Case{
		{
			label:    "Should pass",
			filePath: "./test_shouldPass.json",
			wantErr:  false,
			result: []testcase.TestCase{
				{
					TestLabel: "1",
					ActionUID: "test_action",
					ActionURL: "http://localhost:3000",
					InputFields: map[string]any{
						"label": "Brian Halligan",
						"value": json.Number("1234"),
					},
					ObjectID:              123,
					ObjectType:            "CONTACT",
					PortalID:              111,
					ExpectedExecutionRule: "Success",
				},
				{
					TestLabel: "2",
					ActionUID: "new_action",
					ActionURL: "https://api.hubapi.com/crm",
					InputFields: map[string]any{
						"label": "Maria Johnson",
						"value": json.Number("9876"),
					},
					ObjectID:              987,
					ObjectType:            "CONTACT",
					PortalID:              222,
					ExpectedExecutionRule: "Failure",
				},
			},
		},
		{
			label:    "Error - no file",
			filePath: "./test_noFile.json",
			wantErr:  true,
		},
		{
			label:    "Error - no records",
			filePath: "./test_noRecords.json",
			wantErr:  true,
		},
		{
			label:    "Error - invalid inputs",
			filePath: "./test_invalidInputs.json",
			wantErr:  true,
		},
		{
			label:    "Error - no actionUID",
			filePath: "./test_noActionUID.json",
			wantErr:  true,
		},
		{
			label:    "Error - no testLabel",
			filePath: "./test_noTestLabel.json",
			wantErr:  true,
		},
		{
			label:    "Error - no portalID",
			filePath: "./test_noPortalID.json",
			wantErr:  true,
		},
		{
			label:    "Error - no objectType",
			filePath: "./test_noObjectType.json",
			wantErr:  true,
		},
		{
			label:    "Error - no objectID",
			filePath: "./test_noObjectID.json",
			wantErr:  true,
		},
	}

	for _, thisCase := range cases {
		t.Run(thisCase.label, func(t *testing.T) {
			result, err := testcase.Parse(thisCase.filePath)
			if err != nil != thisCase.wantErr {
				t.Error(err.Error())
				return
			}
			matches := slices.EqualFunc(result, thisCase.result, func(a, b testcase.TestCase) bool {
				if a.ExpectedExecutionRule != b.ExpectedExecutionRule {
					fmt.Println("Mismatched executionLabel")
					return false
				}
				if !maps.Equal(a.InputFields, b.InputFields) {
					fmt.Println("Mismatched InputFields:", a.InputFields, b.InputFields)
					fmt.Println("a")
					for k, v := range a.InputFields {
						fmt.Printf("%s: %v (%T)\n", k, v, v)
					}
					fmt.Println("b")
					for k, v := range b.InputFields {
						fmt.Printf("%s: %v (%T)\n", k, v, v)
					}
					return false
				}
				return true
			})
			if !matches {
				t.Errorf("Expected %v\ngot %v", thisCase.result, result)
			}
		})
	}
}
