package hsserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/jshelper"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func (s *HSServer) createRequest(
	ctx context.Context,
	actionDef actiondefinition.ActionDefinition,
	testCase testcase.TestCase,
) (req *http.Request, callbackId string, err error) {
	method, url, callbackId, headers, serializedBody, err := getRequestData(actionDef, testCase)
	if err != nil {
		return nil, "", err
	}

	req, err = http.NewRequestWithContext(
		ctx,
		method,
		url,
		bytes.NewReader(serializedBody),
	)
	if err != nil {
		return nil, "", err
	}

	if _, ok := headers["content-type"]; !ok {
		req.Header.Add("Content-Type", "application/json")
	}

	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	req.Header.Add("x-hubspot-signature", "v2")
	req.Header.Add(
		"X-HubSpot-Signature",
		signRequestV2(
			s.clientSecret,
			method,
			url,
			serializedBody,
		),
	)
	timestamp := time.Now().UnixMilli()
	v3Signature := signRequestV3(s.clientSecret, method, url, serializedBody, timestamp)
	req.Header.Add("x-hubspot-signature-v3", v3Signature)
	req.Header.Add("x-hubspot-request-timestamp", strconv.Itoa(int(timestamp)))

	if callbackId == "" {
		callbackId = v3Signature
	}

	return req, callbackId, nil
}

func getRequestData(
	actionDef actiondefinition.ActionDefinition,
	testCase testcase.TestCase,
) (method, url, callbackId string, headers map[string]string, serializedBody []byte, err error) {
	type Origin struct {
		PortalID int `json:"portalId"`
	}

	switch test := testCase.Test.(type) {
	case testcase.ActionTest:
		type Object struct {
			ObjectID   int    `json:"objectId"`
			ObjectType string `json:"objectType"`
		}
		callbackId = fmt.Sprintf("ap-%d-%d-0-1", rand.Int(), rand.Int())
		method = "POST"
		url = actionDef.Config.ActionURL
		if test.ActionURL != "" {
			url = test.ActionURL
		}
		headers = map[string]string{}

		origin := Origin{
			PortalID: testCase.PortalID,
		}
		object := Object{
			ObjectID:   test.ObjectID,
			ObjectType: test.ObjectType,
		}

		if preActionFunction := actionDef.GetPreActionFunction(); preActionFunction != nil {
			jsCallback, err := jshelper.RunPreActionFunction(
				jshelper.PreActionEvent{
					WebhookURL:  url,
					CallbackID:  callbackId,
					InputFields: test.InputFields,
					Object:      object,
					Origin:      origin,
				},
				preActionFunction.SourceCode(),
			)
			if err != nil {
				return "", "", "", nil, nil, err
			}

			method = string(jsCallback.HttpMethod)
			url = jsCallback.WebhookURL
			if jsCallback.HttpHeaders != nil {
				headers = jsCallback.HttpHeaders
			}
			headers["content-type"] = jsCallback.ContentType
			headers["accept"] = jsCallback.Accept

			serializedBody, err = json.Marshal(jsCallback.Body)
			if err != nil {
				return "", "", "", nil, nil, err
			}

		} else {
			type Body struct {
				CallbackID string         `json:"callbackId"`
				Fields     map[string]any `json:"fields"`
				Object     Object         `json:"object"`
				Origin     Origin         `json:"origin"`
			}
			body := Body{
				CallbackID: callbackId,
				Fields:     test.InputFields,
				Object:     object,
				Origin:     origin,
			}
			serializedBody, err = json.Marshal(body)
			if err != nil {
				return "", "", "", nil, nil, err
			}

		}
	case testcase.OptionTest:
		method = "POST"
		url = actionDef.Config.ActionURL
		if test.OptionsURL != "" {
			url = test.OptionsURL
		}
		headers = map[string]string{}

		origin := Origin{
			PortalID: testCase.PortalID,
		}

		if preOptionFunction := actionDef.GetPreOptionFunction(test.InputFieldName); preOptionFunction != nil {
			jsCallback, err := jshelper.RunPreOptionFunction(
				jshelper.PreOptionEvent{
					WebhookURL:     url,
					InputFields:    test.InputFields,
					InputFieldName: test.InputFieldName,
					ObjectTypeID:   test.ObjectTypeID,
					Origin:         origin,
				},
				preOptionFunction.SourceCode(),
			)
			if err != nil {
				return "", "", "", nil, nil, err
			}

			method = string(jsCallback.HttpMethod)
			url = jsCallback.WebhookURL
			if jsCallback.HttpHeaders != nil {
				headers = jsCallback.HttpHeaders
			}
			headers["content-type"] = jsCallback.ContentType
			headers["accept"] = jsCallback.Accept

			serializedBody, err = json.Marshal(jsCallback.Body)
			if err != nil {
				return "", "", "", nil, nil, err
			}

		} else {
			type Body struct {
				InputFields    map[string]testcase.OptionInputField `json:"inputFields"`
				InputFieldName string                               `json:"inputFieldName"`
				ObjectTypeID   string                               `json:"objectTypeId"`
				Origin         Origin                               `json:"origin"`
			}
			body := Body{
				InputFields:    test.InputFields,
				InputFieldName: test.InputFieldName,
				ObjectTypeID:   test.ObjectTypeID,
				Origin:         origin,
			}
			serializedBody, err = json.Marshal(body)
			if err != nil {
				return "", "", "", nil, nil, err
			}

		}
		return method, url, "", headers, serializedBody, nil
	default:
		return "", "", "", nil, nil, fmt.Errorf("Unknown test type: %v", testCase.Test)
	}

	return method, url, callbackId, headers, serializedBody, nil
}
