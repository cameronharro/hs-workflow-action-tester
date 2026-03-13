package actiondefinition

type FunctionType string

const (
	PreActionExecution  FunctionType = "PRE_ACTION_EXECUTION"
	PostActionExecution FunctionType = "POST_ACTION_EXECUTION"
	PreFetchOptions     FunctionType = "PRE_FETCH_OPTIONS"
	PostFetchOptions    FunctionType = "POST_FETCH_OPTIONS"
)

type Function interface {
	Type() FunctionType
}

type ActionFunction struct {
	FunctionType   FunctionType `json:"functionType"`
	FunctionSource string       `json:"functionSource"`
}

func (f ActionFunction) Type() FunctionType {
	return f.FunctionType
}

type OptionFunction struct {
	FunctionType   FunctionType `json:"functionType"`
	Id             string       `json:"id"`
	FunctionSource string       `json:"functionSource"`
}

func (f OptionFunction) Type() FunctionType {
	return f.FunctionType
}
