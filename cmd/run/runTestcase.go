package main

import (
	"fmt"
	"slices"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/hsserver"
	"github.com/cameronharro/hs-workflow-tester/internal/testcases"
)

func getDefForCase(testCase testcases.TestCase, actionDefs []actiondefinition.ActionDefinition) (actiondefinition.ActionDefinition, error) {
	defIndex := slices.IndexFunc(actionDefs, func(def actiondefinition.ActionDefinition) bool {
		return def.Uid == testCase.ActionUID
	})
	if defIndex == -1 {
		return actiondefinition.ActionDefinition{}, fmt.Errorf("Could not find action definition uid: %s\n", testCase.ActionUID)
	}
	return actionDefs[defIndex], nil
}

func validateCaseAgainstDef(testCase testcases.TestCase, actionDef actiondefinition.ActionDefinition) error {
	for _, actionInput := range actionDef.Config.InputFields {
		if testCaseMissingRequiredInput(testCase, actionInput) {
			return fmt.Errorf("Invalid testCase %v: Missing required field %v\n", testCase, actionInput)
		}
	}
	return nil
}

func testCaseMissingRequiredInput(testCase testcases.TestCase, actionInput actiondefinition.InputField) bool {
	return actionInput.IsRequired && !slices.ContainsFunc(testCase.InputFields, func(testInput testcases.InputField) bool {
		return testInput.Name == actionInput.TypeDefinition.GetName() && testInput.Value != ""
	})
}

func runTestCase(
	server *hsserver.HSServer,
	testCase testcases.TestCase,
	actionDefs []actiondefinition.ActionDefinition,
) error {
	actionDef, err := getDefForCase(testCase, actionDefs)
	if err != nil {
		return err
	}

	if err = validateCaseAgainstDef(testCase, actionDef); err != nil {
		return err
	}

	server.SendRequest(actionDef, testCase)
	return nil
}
