package jshelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type Envelope struct {
	Event    Event  `json:"event"`
	Function string `json:"function"`
}

type FunctionType string

const (
	PreActionExecution  FunctionType = "PRE_ACTION_EXECUTION"
	PostActionExecution FunctionType = "POST_ACTION_EXECUTION"
	PreFetchOptions     FunctionType = "PRE_FETCH_OPTIONS"
	PostFetchOptions    FunctionType = "POST_FETCH_OPTIONS"
)

type Event interface {
	getEventType() FunctionType
}

type CallbackData interface {
	getCallbackType() FunctionType
}

type Method string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Patch  Method = "PATCH"
	Put    Method = "PUT"
	Delete Method = "DELETE"
)

func getAllMethods() []Method {
	return []Method{Get, Post, Patch, Put, Delete}
}

type RequestParams struct {
	URL     string            `json:"url"`
	Method  Method            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    map[string]any    `json:"body"`
}

func RunFunction(event Event, function string) (CallbackData, error) {
	switch event.(type) {
	case PreActionEvent:
		return spawn(event, function, validatePreAction)

	default:
		return PreActionCallback{}, fmt.Errorf("Invalid event type for %v", event)
	}
}

func spawn[T Event, V CallbackData](event T, function string, validator func(d []byte) (V, error)) (callbackData V, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	c := exec.CommandContext(ctx, "deno", "run", "--deny-read", "--deny-write", "--deny-net", "--deny-env", "--deny-run", "--deny-ffi", "--deny-sys", "--deny-import", "./jsHelper.ts")
	envelope := Envelope{
		Event:    event,
		Function: function,
	}
	jsonBytes, err := json.Marshal(envelope)
	if err != nil {
		return *new(V), err
	}

	c.Stdin = bytes.NewReader(jsonBytes)
	out, err := c.Output()
	if err != nil {
		return *new(V), err
	}

	callbackData, err = validator(out)
	if err != nil {
		return *new(V), err
	}

	return callbackData, nil
}
