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
	Data RequestParams `json:"data"`
	Code string        `json:"code"`
}

type Method string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Patch  Method = "PATCH"
	Put    Method = "PUT"
	Delete Method = "DELETE"
)

type RequestParams struct {
	URL     string            `json:"url"`
	Method  Method            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    map[string]any    `json:"body"`
}

func Spawn(data RequestParams, code string) (transformedData RequestParams, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	c := exec.CommandContext(ctx, "deno", "run", "--deny-read", "--deny-write", "--deny-net", "--deny-env", "--deny-run", "--deny-ffi", "--deny-sys", "--deny-import", "./jsHelper.ts")
	envelope := Envelope{
		Data: data,
		Code: code,
	}
	jsonBytes, err := json.Marshal(envelope)
	if err != nil {
		return RequestParams{}, err
	}

	c.Stdin = bytes.NewReader(jsonBytes)
	out, err := c.Output()
	if err != nil {
		return RequestParams{}, err
	}

	fmt.Println(string(out))
	return RequestParams{}, nil
}
