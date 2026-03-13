package hsserver

type TextTypeDefintion struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	FieldType string `json:"fieldType"`
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
type NumberTypeDefinitino struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type InputField struct {
	TypeDefinition any  `json:"typeDefinition"`
	IsRequired     bool `json:"isRequired"`
}

type FunctionType string

const (
	PreActionExecution  FunctionType = "PRE_FUNCTION_EXECUTION"
	PostActionExecution FunctionType = "POST_FUNCTION_EXECUTION"
	PreFetchOptions     FunctionType = "PRE_FETCH_OPTIONS"
	POSTFetchOptions    FunctionType = "POST_FETCH_OPTIONS"
)

type ActionFunction struct {
	FunctionType   FunctionType `json:"functionType"`
	FunctionSource string       `json:"functionSource"`
}
type OptionFunction struct {
	FunctionType   FunctionType `json:"functionType"`
	Id             string       `json:"id"`
	FunctionSource string       `json:"functionSource"`
}
type ActionConfig struct {
	InputFields []any `json:"inputFields"`
	Functions   []any `json:"functions"`
}
type ActionDefinition struct {
	Uid    string       `json:"uid"`
	Type   string       `json:"type"`
	Config ActionConfig `json:"config"`
}
