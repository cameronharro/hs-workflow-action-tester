package actiondefinition

type ExecutionRule struct {
	LabelName  string         `json:"labelName"`
	Conditions map[string]any `json:"conditions"`
}

func (r ExecutionRule) matchesOutput(
	outputFields map[string]any,
) bool {
RuleEvaluation:
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
			if concreteExVal != concreteOutputVal {
				return false
			}
		case []any:
			for _, concreteExValElement := range concreteExVal {
				elementStr, ok := concreteExValElement.(string)
				if !ok {
					continue
				}

				if elementStr == concreteOutputVal {
					continue RuleEvaluation
				}
			}
			return false
		default:
			return false
		}
	}
	return true
}

func (d ActionDefinition) GetMatchingExecutionRule(outputFields map[string]any) string {
	if outputFields == nil {
		return ""
	}

	for _, executionRule := range d.Config.ExecutionRules {
		if executionRule.matchesOutput(outputFields) {
			return executionRule.LabelName
		}
	}
	return ""
}
