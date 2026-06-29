package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cameronharro/hs-workflow-tester/internal/hsserver"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func main() {
	flags := initFlags()

	testCases, err := testcase.Parse(flags.TestCasePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	actionDefinitions, err := getActionDefs(flags.HSProjectRoot)
	if err != nil {
		log.Fatal(err.Error())
	}

	server := hsserver.NewHSServer(flags.ClientSecret, flags.AsyncListenerPort, flags.Timeout)

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
