package hsserver

import (
	"encoding/json"
	"errors"
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
	var resultErr error
	switch test := testCase.Test.(type) {
	case testcase.ActionTest:

		if err := testCaseMatchesObjectType(test, actionDef); err != nil {
			resultErr = errors.Join(resultErr, err)
		}

		for _, actionInput := range actionDef.Config.InputFields {
			if err := testCaseHasRequiredInput(test, actionInput); err != nil {
				resultErr = errors.Join(resultErr, err)
			}
		}

		if err := testCaseInputsMatchTypes(test, actionDef); err != nil {
			resultErr = errors.Join(resultErr, err)
		}
	case testcase.OptionTest:
		return fmt.Errorf("Option test validation not implemented yet")
	default:
		return fmt.Errorf("unknown testcase type for validation")
	}

	return resultErr
}

func testCaseHasRequiredInput(
	actionTest testcase.ActionTest,
	actionInput actiondefinition.InputField,
) error {
	_, testCaseHasInput := actionTest.InputFields[actionInput.TypeDefinition.GetName()]
	if !testCaseHasInput && actionInput.IsRequired {
		return fmt.Errorf("Missing required field %s", actionInput.TypeDefinition.GetName())
	}
	return nil
}

func testCaseMatchesObjectType(
	actionTest testcase.ActionTest,
	actionDef actiondefinition.ActionDefinition,
) error {
	allowsAllObjects := len(actionDef.Config.ObjectTypes) == 0
	if allowsAllObjects {
		return nil
	}

	if slices.Contains(actionDef.Config.ObjectTypes, actionTest.ObjectType) {
		return nil
	}

	return fmt.Errorf("objectType %s is not in configured permitted objects: %v", actionTest.ObjectType, actionDef.Config.ObjectTypes)
}

func testCaseInputsMatchTypes(
	actionTest testcase.ActionTest,
	actionDef actiondefinition.ActionDefinition,
) error {
	var resultErr error
	for key, value := range actionTest.InputFields {
		actionInputIdx := slices.IndexFunc(actionDef.Config.InputFields, func(field actiondefinition.InputField) bool {
			return field.TypeDefinition.GetName() == key
		})
		if actionInputIdx == -1 {
			resultErr = errors.Join(resultErr, fmt.Errorf("provided inputField %s is not an input option", key))
			continue
		}

		compareTypes := func(actionInput actiondefinition.InputField) error {
			switch typeDef := actionInput.TypeDefinition.(type) {
			case actiondefinition.EnumTypeDefinition:
				s, ok := value.(string)
				if !ok {
					return fmt.Errorf("%s [enum]: %v is not a string ", key, value)
				}
				if typeDef.ExternalOptions || typeDef.OptionsURL != "" {
					break
				}
				if !slices.ContainsFunc(typeDef.Options, func(opt actiondefinition.Option) bool {
					return opt.Value == s
				}) {
					return fmt.Errorf("%s [enum]: %s is not an option: %v", key, s, typeDef.Options)
				}
			case actiondefinition.NumberTypeDefinition:
				if _, ok := value.(json.Number); !ok {
					return fmt.Errorf("%s [number]: %v is not a number]", key, value)
				}
			case actiondefinition.StringTypeDefinition:
				if _, ok := value.(string); !ok {
					return fmt.Errorf("%s [string]: %v is not a string", key, value)
				}
			default:
				return fmt.Errorf("%v is an unknown actionInput type definition", typeDef)
			}
			return nil
		}

		if err := compareTypes(actionDef.Config.InputFields[actionInputIdx]); err != nil {
			resultErr = errors.Join(resultErr, err)
		}
	}

	return resultErr
}
