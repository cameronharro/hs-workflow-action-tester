package testcase_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
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
					PortalID:  111,
					TestType:  testcase.Action,
					Test: testcase.ActionTest{

						ActionURL: "http://localhost:3000",
						InputFields: map[string]any{
							"label": "Brian Halligan",
							"value": json.Number("1234"),
						},
						ObjectID:              123,
						ObjectType:            "CONTACT",
						ExpectedExecutionRule: "Success",
					},
				},
				{
					TestLabel: "2",
					ActionUID: "new_action",
					PortalID:  222,
					TestType:  testcase.Action,
					Test: testcase.ActionTest{
						ActionURL: "https://api.hubapi.com/crm",
						InputFields: map[string]any{
							"label": "Maria Johnson",
							"value": json.Number("9876"),
						},
						ObjectID:              987,
						ObjectType:            "CONTACT",
						ExpectedExecutionRule: "Failure",
					},
				},
				{
					TestLabel: "options",
					ActionUID: "old_action",
					PortalID:  222,
					TestType:  testcase.Option,
					Test: testcase.OptionTest{
						OptionsURL:     "https://api.hubapi.com/crm",
						InputFieldName: "custom_enum",
						InputFields: map[string]testcase.OptionInputField{
							"name": testcase.ObjectPropertyInputField{
								FieldType:    testcase.ObjectProperty,
								PropertyName: "firstname",
							},
							"status": testcase.StaticValueInputField{
								FieldType: testcase.StaticValue,
								Value:     "active",
							},
						},
						ObjectTypeID: "0-1",
					},
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
			for i, targetCase := range thisCase.result {
				var thisErr error
				if len(result) <= i {
					thisErr = fmt.Errorf("testCase %s missing counterpart in results: %v", targetCase.TestLabel, result)
					err = errors.Join(thisErr)
					break
				}

				resultCase := result[i]
				if targetCase.ActionUID != resultCase.ActionUID {
					thisErr = mismatchedErr(thisErr, "actionUID", targetCase.ActionUID, resultCase.ActionUID)
				}

				switch test := targetCase.Test.(type) {
				case testcase.OptionTest:
					thisErr = errors.Join(thisErr, optTestEq(test, resultCase.Test))
				case testcase.ActionTest:
					thisErr = errors.Join(thisErr, actionTestEq(test, resultCase.Test))
				default:
					thisErr = errors.Join(fmt.Errorf("unknown test type"))
				}

				if thisErr != nil {
					thisErr = fmt.Errorf("testCase %s: %w", targetCase.TestLabel, thisErr)
					err = errors.Join(err, thisErr)
				}
			}
			if err != nil != thisCase.wantErr {
				t.Error(err.Error())
			}
		})
	}
}

func optTestEq(target testcase.OptionTest, result any) error {
	resultTest, ok := result.(testcase.OptionTest)
	if !ok {
		return fmt.Errorf("not OptionTest: %v", result)
	}

	var err error
	if target.InputFieldName != resultTest.InputFieldName {
		err = mismatchedErr(err, "inputFieldName", target.InputFieldName, resultTest.InputFieldName)
	}

	if target.ObjectTypeID != resultTest.ObjectTypeID {
		err = mismatchedErr(err, "objectTypeId", target.ObjectTypeID, resultTest.ObjectTypeID)
	}

	if target.OptionsURL != resultTest.OptionsURL {
		err = mismatchedErr(err, "optionsURL", target.OptionsURL, resultTest.OptionsURL)
	}

	if !maps.Equal(target.InputFields, resultTest.InputFields) {
		err = mismatchedErr(err, "inputFields", target.InputFields, resultTest.InputFields)
	}
	return err
}

func actionTestEq(target testcase.ActionTest, result any) error {
	resultTest, ok := result.(testcase.ActionTest)
	if !ok {
		return fmt.Errorf("not ActionTest: %v", result)
	}

	var err error
	if target.ExpectedExecutionRule != resultTest.ExpectedExecutionRule {
		err = mismatchedErr(err, "expectedExecutionRule", target.ExpectedExecutionRule, resultTest.ExpectedExecutionRule)
	}

	if target.ObjectID != resultTest.ObjectID {
		err = mismatchedErr(err, "objectID", target.ObjectID, resultTest.ObjectID)
	}

	if target.ObjectType != resultTest.ObjectType {
		err = mismatchedErr(err, "objectType", target.ObjectType, resultTest.ObjectType)
	}

	if target.ActionURL != resultTest.ActionURL {
		err = mismatchedErr(err, "actionURL", target.ActionURL, resultTest.ActionURL)
	}

	if !maps.Equal(target.InputFields, resultTest.InputFields) {
		err = mismatchedErr(err, "inputFields", target.InputFields, resultTest.InputFields)
	}
	return err
}

func mismatchedErr(err error, key string, expectedVal, receivedVal any) error {
	return errors.Join(err, fmt.Errorf(
		"Mismatched %s: expected %v, got %v",
		key,
		expectedVal,
		receivedVal,
	))
}
