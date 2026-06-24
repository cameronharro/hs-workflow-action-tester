package actiondefinition

import "slices"

type FunctionType string

const (
	PreActionExecution  FunctionType = "PRE_ACTION_EXECUTION"
	PostActionExecution FunctionType = "POST_ACTION_EXECUTION"
	PreFetchOptions     FunctionType = "PRE_FETCH_OPTIONS"
	PostFetchOptions    FunctionType = "POST_FETCH_OPTIONS"
)

type Function interface {
	Type() FunctionType
	SourceCode() string
}

func (d ActionDefinition) getActionFunction(t FunctionType) Function {
	if !slices.Contains([]FunctionType{PreActionExecution, PostActionExecution}, t) {
		return nil
	}

	index := slices.IndexFunc(d.Config.Functions, func(ele Function) bool {
		return ele.Type() == PreActionExecution
	})
	if index == -1 {
		return nil
	}
	return d.Config.Functions[index]
}

func (d ActionDefinition) GetPreActionFunction() Function {
	return d.getActionFunction(PreActionExecution)
}

func (d ActionDefinition) GetPostActionFunction() Function {
	return d.getActionFunction(PostActionExecution)
}

type ActionFunction struct {
	FunctionType   FunctionType `json:"functionType"`
	FunctionSource string       `json:"functionSource"`
}

func (f ActionFunction) Type() FunctionType {
	return f.FunctionType
}

func (f ActionFunction) SourceCode() string {
	return f.FunctionSource
}

type OptionFunction struct {
	FunctionType   FunctionType `json:"functionType"`
	Id             string       `json:"id"`
	FunctionSource string       `json:"functionSource"`
}

func (f OptionFunction) Type() FunctionType {
	return f.FunctionType
}

func (f OptionFunction) SourceCode() string {
	return f.FunctionSource
}
