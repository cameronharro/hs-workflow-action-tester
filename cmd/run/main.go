package main

import (
	"fmt"
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

	for _, testCase := range testCases {
		err = runTestCase(server, testCase, actionDefinitions)
		if err != nil {
			fmt.Print(err)
		}
	}
}
