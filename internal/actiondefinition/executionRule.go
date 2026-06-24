package actiondefinition

import (
	"slices"
)

type ExecutionRule struct {
	LabelName  string         `json:"labelName"`
	Conditions map[string]any `json:"conditions"`
}

func (r ExecutionRule) MatchesOutput(
	outputFields map[string]any,
) bool {
	for executionKey, executionValue := range r.Conditions {
		outputValue, ok := outputFields[executionKey]
		if !ok {
			return false
		}

		concreteOutputVal, ok := outputValue.(string)
		if !ok {
			return false
		}

		switch concreteExVal := executionValue.(type) {
		case string:
			if concreteExVal != outputValue {
				return false
			}
		case []string:
			if !slices.Contains(concreteExVal, concreteOutputVal) {
				return false
			}
		default:
			return false
		}
	}
	return true
}
