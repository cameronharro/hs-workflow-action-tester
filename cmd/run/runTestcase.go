package main

import (
	"fmt"
	"slices"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/hsserver"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func getDefForCase(
	testCase testcase.TestCase,
	actionDefs []actiondefinition.ActionDefinition,
) (actiondefinition.ActionDefinition, error) {
	defIndex := slices.IndexFunc(actionDefs, func(def actiondefinition.ActionDefinition) bool {
		return def.Uid == testCase.ActionUID
	})
	if defIndex == -1 {
		err := fmt.Errorf("Could not find action definition uid: %s\n", testCase.ActionUID)
		return actiondefinition.ActionDefinition{}, err
	}
	return actionDefs[defIndex], nil
}

func validateCaseAgainstDef(
	testCase testcase.TestCase,
	actionDef actiondefinition.ActionDefinition,
) error {
	for _, actionInput := range actionDef.Config.InputFields {
		if testCaseMissingRequiredInput(testCase, actionInput) {
			return fmt.Errorf("Invalid testCase %v: Missing required field %v\n", testCase, actionInput)
		}
		if !slices.Contains(actionDef.Config.ObjectTypes, testCase.ObjectType) {
			return fmt.Errorf("Invalid testCase %v: objectType %s not in config permitted objects: %v", testCase, testCase.ObjectType, actionDef.Config.ObjectTypes)
		}
	}
	return nil
}

func testCaseMissingRequiredInput(
	testCase testcase.TestCase,
	actionInput actiondefinition.InputField,
) bool {
	_, testCaseHasInput := testCase.InputFields[actionInput.TypeDefinition.GetName()]
	return testCaseHasInput || !actionInput.IsRequired
}

func runTestCase(
	server *hsserver.HSServer,
	testCase testcase.TestCase,
	actionDefs []actiondefinition.ActionDefinition,
) error {
	actionDef, err := getDefForCase(testCase, actionDefs)
	if err != nil {
		return err
	}

	if err = validateCaseAgainstDef(testCase, actionDef); err != nil {
		return err
	}

	_, _, err = server.CreateRequest(actionDef, testCase)
	if err != nil {
		return err
	}

	return nil
}
