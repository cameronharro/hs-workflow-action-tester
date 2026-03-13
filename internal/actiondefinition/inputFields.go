package actiondefinition

import (
	"encoding/json"
	"errors"
)

type FieldTypeDefinition interface {
	GetName() string
}

type StringTypeDefinition struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	FieldType string `json:"fieldType"`
}

func (f StringTypeDefinition) GetName() string {
	return f.Name
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
type EnumTypeDefinition struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	FieldType string   `json:"fieldType"`
	Options   []Option `json:"options"`
}

func (e EnumTypeDefinition) GetName() string {
	return e.Name
}

type NumberTypeDefinition struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (n NumberTypeDefinition) GetName() string {
	return n.Name
}

type InputField struct {
	TypeDefinition FieldTypeDefinition `json:"typeDefinition"`
	IsRequired     bool                `json:"isRequired"`
}

func (c *InputField) UnmarshalJSON(data []byte) error {
	var raw struct {
		TypeDefinition json.RawMessage `json:"typeDefinition"`
		IsRequired     bool            `json:"isRequired"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.IsRequired = raw.IsRequired

	var disc struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw.TypeDefinition, &disc); err != nil {
		return err
	}

	switch disc.Type {
	case "string":
		var def StringTypeDefinition
		if err := json.Unmarshal(raw.TypeDefinition, &def); err != nil {
			return err
		}
		c.TypeDefinition = def
	case "number":
		var def NumberTypeDefinition
		if err := json.Unmarshal(raw.TypeDefinition, &def); err != nil {
			return err
		}
		c.TypeDefinition = def
	case "enumeration":
		var def EnumTypeDefinition
		if err := json.Unmarshal(raw.TypeDefinition, &def); err != nil {
			return err
		}
		c.TypeDefinition = def
	default:
		return errors.New("invalid InputField Type")
	}
	return nil
}
