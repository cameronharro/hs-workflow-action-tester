package testcase

import (
	"encoding/json"
	"errors"
	"fmt"
)

type OptionTest struct {
	ObjectTypeID    string `json:"objectTypeId"`
	OptionsURL      string `json:"optionsURL"`
	InputFieldName  string `json:"inputFieldName"`
	InputFields     map[string]OptionInputField
	ExpectedOptions []Option `json:"expectedOptions"`
}

type Option struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

func (t OptionTest) Validate() error {
	var result error
	if t.ObjectTypeID == "" {
		result = errors.Join(result, errors.New("Missing objectTypeId"))

	}

	if t.InputFieldName == "" {
		result = errors.Join(result, errors.New("Missing inputFieldName"))
	}

	for k, f := range t.InputFields {
		if err := f.ValidateField(); err != nil {
			result = errors.Join(result, fmt.Errorf("%s: %w", k, err))
		}
	}

	return result
}

func (t *OptionTest) UnmarshalJSON(data []byte) error {
	type Peek struct {
		ObjectTypeID    string `json:"objectTypeId"`
		OptionsURL      string `json:"optionsURL"`
		InputFieldName  string `json:"inputFieldName"`
		InputFields     map[string]json.RawMessage
		ExpectedOptions []Option `json:"expectedOptions"`
	}
	peek := Peek{}
	err := json.Unmarshal(data, &peek)
	if err != nil {
		return err
	}

	t.ObjectTypeID = peek.ObjectTypeID
	t.OptionsURL = peek.OptionsURL
	t.InputFieldName = peek.InputFieldName
	t.InputFields = map[string]OptionInputField{}
	t.ExpectedOptions = peek.ExpectedOptions
	for k, v := range peek.InputFields {
		type Peek struct {
			FieldType OptionInputFieldType `json:"type"`
		}
		peek := Peek{}
		err := json.Unmarshal(v, &peek)
		if err != nil {
			return err
		}

		var optionInputField OptionInputField
		switch peek.FieldType {
		case ObjectProperty:
			objPropField := ObjectPropertyInputField{}
			err = json.Unmarshal(v, &objPropField)
			if err != nil {
				return err
			}
			optionInputField = objPropField
		case StaticValue:
			staticValField := StaticValueInputField{}
			err = json.Unmarshal(v, &staticValField)
			if err != nil {
				return err
			}
			optionInputField = staticValField
		default:
			return fmt.Errorf("Unknown Option Input Field Type: %s", peek.FieldType)
		}

		t.InputFields[k] = optionInputField
	}
	return nil
}

type OptionInputField interface {
	ValidateField() error
}

type StaticValueInputField struct {
	FieldType OptionInputFieldType `json:"type"`
	Value     string               `json:"value"`
}

func (f StaticValueInputField) ValidateField() error {
	return nil
}

type ObjectPropertyInputField struct {
	FieldType    OptionInputFieldType `json:"type"`
	PropertyName string               `json:"propertyName"`
}

type OptionInputFieldType string

const (
	ObjectProperty OptionInputFieldType = "OBJECT_PROPERTY"
	StaticValue    OptionInputFieldType = "STATIC_VALUE"
)

func (f ObjectPropertyInputField) ValidateField() error {
	if f.PropertyName == "" {
		return errors.New("Missing propertyName")
	}
	return nil
}
