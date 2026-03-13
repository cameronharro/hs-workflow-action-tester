package actiondefinition

import (
	"encoding/json"
)

func Parse(data []byte) (*ActionDefinition, error) {
	result := ActionDefinition{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
