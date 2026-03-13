package hsserver

import (
	"fmt"
	"testing"
)

func TestSignRequest(t *testing.T) {
	type TestCase struct {
		clientSecret string
		method       string
		uri          string
		body         string
		result       string
	}

	testCases := []TestCase{
		{
			clientSecret: "fa768533-04ea-4d9b-becf-38cef1f96f0f",
			method:       "POST",
			uri:          "https://webhook.site/c5db800b-e3c3-4a16-af81-1e19f6d8b6a0",
			body:         `{"callbackId":"ap-50467010-2384423471376-1-0","origin":{"portalId":50467010,"userId":null,"actionDefinitionId":259551814,"actionDefinitionVersion":1,"actionExecutionIndexIdentifier":{"enrollmentId":2384423471376,"actionExecutionIndex":0},"extensionDefinitionVersionId":1,"extensionDefinitionId":259551814},"context":{"workflowId":1791092093,"actionId":1,"actionExecutionIndexIdentifier":{"enrollmentId":2384423471376,"actionExecutionIndex":0},"source":"WORKFLOWS"},"object":{"objectId":153847279758,"objectType":"CONTACT"},"fields":{"label":"Brian Halligan (Sample Contact)","value":"153847279758"},"inputFields":{"label":"Brian Halligan (Sample Contact)","value":"153847279758"},"typedInputs":{"label":{"value":"Brian Halligan (Sample Contact)","type":"STRING"},"value":{"value":"153847279758","type":"STRING"}}}`,
			result:       "d0c157dec62681b3ef5b9f7a1d66a9a63cfc67984b351a0f9d736e82ebdb0f03",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test %d:\n", i+1), func(t *testing.T) {
			result := signRequestV2(testCase.clientSecret, testCase.method, testCase.uri, testCase.body)
			if result != testCase.result {
				t.Errorf("Expected %s, got %s", testCase.result, result)
			}
		})
	}
}
