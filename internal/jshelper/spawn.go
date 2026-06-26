package jshelper

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

//go:embed jsHelper.cts
var jsHelper []byte

func spawn[T Event, V CallbackData](
	event T,
	function string,
	validator func(d []byte) (V, error),
) (callbackData V, err error) {
	errorOut := func(msg string) (V, error) {
		return *new(V), errors.New("jshelper Error: " + msg)
	}

	if !strings.Contains(function, "exports.main = ") {
		return errorOut(fmt.Sprintf("%s definition lacks exports.main assigment: %s\n", event.getEventType(), function))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tmp, err := os.CreateTemp("", "embedded-jsHelper-*.cts")
	if err != nil {
		errorOut(err.Error())
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(jsHelper); err != nil {
		_ = tmp.Close()
		errorOut(err.Error())
	}

	if err := tmp.Close(); err != nil {
		errorOut(err.Error())
	}

	c := exec.CommandContext(
		ctx,
		"deno",
		"run",
		"--allow-read",
		"--deny-write",
		"--deny-net",
		"--deny-env",
		"--deny-run",
		"--deny-ffi",
		"--deny-sys",
		"--deny-import",
		tmp.Name(),
	)
	envelope := Envelope{
		Event:    event,
		Function: function,
	}
	jsonBytes, err := json.Marshal(envelope)
	if err != nil {
		return errorOut(err.Error())
	}

	c.Stdin = bytes.NewReader(jsonBytes)
	out, err := c.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errorOut("deadline exceeded")
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) || len(exitErr.Stderr) > 0 {
			return errorOut(string(exitErr.Stderr))
		}
		return errorOut(err.Error())
	}

	callbackData, err = validator(out)
	if err != nil {
		return errorOut(err.Error())
	}

	return callbackData, nil
}
