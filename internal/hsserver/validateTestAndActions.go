package hsserver

import (
	"fmt"
	"slices"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
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
		if !testCaseHasRequiredInput(testCase, actionInput) {
			return fmt.Errorf("Invalid testCase %s: Missing required field %s\n", testCase.TestLabel, actionInput.TypeDefinition.GetName())
		}
		allowsAllObjects := len(actionDef.Config.ObjectTypes) == 0
		if !allowsAllObjects && !slices.Contains(actionDef.Config.ObjectTypes, testCase.ObjectType) {
			return fmt.Errorf("Invalid testCase %s: objectType %s not in config permitted objects: %v\n", testCase.TestLabel, testCase.ObjectType, actionDef.Config.ObjectTypes)
		}
	}
	return nil
}

func testCaseHasRequiredInput(
	testCase testcase.TestCase,
	actionInput actiondefinition.InputField,
) bool {
	_, testCaseHasInput := testCase.InputFields[actionInput.TypeDefinition.GetName()]
	return testCaseHasInput || !actionInput.IsRequired
}
