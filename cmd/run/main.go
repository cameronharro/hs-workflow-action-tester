package main

import (
	"log"

	"github.com/cameronharro/hs-workflow-tester/internal/hsserver"
	"github.com/cameronharro/hs-workflow-tester/internal/testcases"
)

func main() {
	testCases, err := testcases.Parse("./testCases.csv")
	if err != nil {
		log.Fatal(err.Error())
	}

	actionDefinitions, err := getActionDefs(".")
	if err != nil {
		log.Fatal(err.Error())
	}

	server := hsserver.NewHSServer("1234", "http://localhost:8080")
	err = server.SendRequest(actionDefinitions[0], testCases[0])
	if err != nil {
		log.Fatal(err.Error())
	}
}
