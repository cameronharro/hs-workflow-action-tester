package jshelper

import "encoding/json"

type InputField interface {
	GetType() string
}

type InputFieldMap map[string]InputField

type ObjectProperty struct {
	Type         string `json:"type"`
	PropertyName string `json:"propertyName"`
}

func (o ObjectProperty) GetType() string {
	return o.Type
}

type StaticValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (v StaticValue) GetType() string {
	return v.Type
}

type FetchOptions struct {
	Q     string `json:"q"`
	After string `json:"after"`
}

func (o FetchOptions) GetType() string {
	return "FETCH_OPTIONS"
}

type PreFetchEvent struct {
	Origin struct {
		PortalID                int `json:"portalId"`
		ActionDefinitionID      int `json:"actionDefinitionId"`
		ActionDefinitionVersion int `json:"actionDefinitionVersion"`
	} `json:"origin"`
	InputFieldName string                `json:"inputFieldName"`
	WebhookURL     string                `json:"webhookUrl"`
	InputFields    map[string]InputField `json:"inputFields"`
}

func (m *InputFieldMap) UnmarshalJSON(data []byte) error {
	var raw struct {
		InputFields map[string]json.RawMessage `json:"inputFields"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw.InputFields {
		var disc struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(v, &disc); err != nil {
			return err
		}
		var inputField InputField
		switch disc.Type {
		case "OBJECT_PROPERTY":
			inputField = ObjectProperty{}
		case "STATIC_VALUE":
			inputField = StaticValue{}
		case "":
			inputField = FetchOptions{}
		}
		if err := json.Unmarshal(v, &inputField); err != nil {
			return err
		}
		(*m)[k] = inputField
	}

	return nil
}

type PreFetchCallback struct {
	WebhookURL  string            `json:"webhookUrl"`
	Body        string            `json:"body"`
	HTTPHeaders map[string]string `json:"httpHeaders"`
	ContentType string            `json:"contentType"`
	Accept      string            `json:"accept"`
	HTTPMethod  string            `json:"httpMethod"`
}
