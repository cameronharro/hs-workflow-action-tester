package main

import (
	"flag"
	"os"
	"time"
)

type Flags struct {
	TestCasePath      string
	ClientSecret      string
	HSProjectRoot     string
	AsyncListenerPort int
	Timeout           time.Duration
}

func initFlags() Flags {
	const required = "[required]"
	testCasePath := flag.String(
		"cases",
		required,
		"The path to the configured test cases",
	)
	clientSecret := flag.String(
		"clientSecret",
		required,
		"A secret used to imitate HubSpot's auth signature. Must match your application's expected value",
	)
	hsProjectRoot := flag.String(
		"project",
		".",
		"The path to the HubSpot project's root",
	)
	asyncListenerPort := flag.Int(
		"port",
		8000,
		"The port where the server will listen for asynchronous callbacks from your application",
	)
	queueTimeout := flag.Duration(
		"timeout",
		5*time.Second,
		"The amount of time the server will wait to receive responses from your test cases",
	)
	flag.Parse()
	if *testCasePath == required || *clientSecret == required {
		flag.Usage()
		os.Exit(1)
	}

	result := Flags{
		TestCasePath:      *testCasePath,
		ClientSecret:      *clientSecret,
		HSProjectRoot:     *hsProjectRoot,
		AsyncListenerPort: *asyncListenerPort,
		Timeout:           *queueTimeout,
	}
	return result
}
