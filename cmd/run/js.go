package main

import (
	"os/exec"
	"strings"
)

func runJS(js string) error {
	cmd := exec.Command("deno", "run", "--no-prompt", "--deny-all", "-")
	cmd.Stdin = strings.NewReader(js)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
