package actiondefinition

import (
	"encoding/json"
	"fmt"
)

func (c *ActionConfig) UnmarshalJSON(data []byte) error {
	var raw struct {
		InputFields []InputField      `json:"inputFields"`
		Functions   []json.RawMessage `json:"functions"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.InputFields = raw.InputFields

	for _, rawFunc := range raw.Functions {
		var disc struct {
			FunctionType FunctionType `json:"functionType"`
		}
		if err := json.Unmarshal(rawFunc, &disc); err != nil {
			return err
		}
		switch disc.FunctionType {
		case PreActionExecution:
			fallthrough
		case PostActionExecution:
			var def ActionFunction
			if err := json.Unmarshal(rawFunc, &def); err != nil {
				return err
			}
			c.Functions = append(c.Functions, def)
		case PreFetchOptions:
			fallthrough
		case PostFetchOptions:
			var def OptionFunction
			if err := json.Unmarshal(rawFunc, &def); err != nil {
				return err
			}
			c.Functions = append(c.Functions, def)
		default:
			return fmt.Errorf("Invalid FunctionType: %q", disc.FunctionType)
		}
	}
	return nil
}

type ActionConfig struct {
	InputFields []InputField `json:"inputFields"`
	Functions   []Function   `json:"functions"`
}
type ActionDefinition struct {
	Uid    string       `json:"uid"`
	Type   string       `json:"type"`
	Config ActionConfig `json:"config"`
}
