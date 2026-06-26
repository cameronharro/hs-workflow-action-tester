package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cameronharro/hs-workflow-tester/internal/hsserver"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func main() {
	testCases, err := testcase.Parse("./testCases.csv")
	if err != nil {
		log.Fatal(err.Error())
	}

	actionDefinitions, err := getActionDefs(".")
	if err != nil {
		log.Fatal(err.Error())
	}

	server := hsserver.NewHSServer("1234", 8080)

	for _, testCase := range testCases {
		go func() {
			err = server.RunTestCase(testCase, actionDefinitions)
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	result := <-server.ResultChan
	if result != nil {
		fmt.Println(result.Error())
		os.Exit(1)
	}
	fmt.Println("All test cases pass!")
	os.Exit(0)
}
