package hsserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/jshelper"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func (s *HSServer) createRequest(
	ctx context.Context,
	actionDef actiondefinition.ActionDefinition,
	testCase testcase.TestCase,
) (req *http.Request, callbackId string, err error) {
	type Origin struct {
		PortalID int `json:"portalId"`
	}
	type Object struct {
		ObjectID   int    `json:"objectId"`
		ObjectType string `json:"objectType"`
	}

	callbackId = fmt.Sprintf("ap-%d-%d-0-1", rand.Int(), rand.Int())

	url := actionDef.Config.ActionURL
	if testCase.ActionURL != "" {
		url = testCase.ActionURL
	}
	origin := Origin{
		PortalID: testCase.PortalID,
	}
	object := Object{
		ObjectID:   testCase.ObjectID,
		ObjectType: testCase.ObjectType,
	}
	if preActionFunction := actionDef.GetPreActionFunction(); preActionFunction != nil {
		jsCallback, err := jshelper.RunPreActionFunction(
			jshelper.PreActionEvent{
				WebhookURL:  url,
				CallbackID:  callbackId,
				InputFields: testCase.InputFields,
				Object:      object,
				Origin:      origin,
			},
			preActionFunction.SourceCode(),
		)
		if err != nil {
			return nil, "", err
		}

		callbackBody, err := json.Marshal(jsCallback.Body)
		if err != nil {
			return nil, "", err
		}

		req, err = http.NewRequestWithContext(
			ctx,
			string(jsCallback.HttpMethod),
			jsCallback.WebhookURL,
			bytes.NewBuffer(callbackBody),
		)
		if err != nil {
			return nil, "", err
		}

		for headerKey, headerValue := range jsCallback.HttpHeaders {
			req.Header.Add(headerKey, headerValue)
		}
		req.Header.Add(
			"X-HubSpot-Signature",
			signRequestV2(
				s.clientSecret,
				string(jsCallback.HttpMethod),
				jsCallback.WebhookURL,
				callbackBody,
			),
		)

	} else {
		type Body struct {
			CallbackID string         `json:"callbackId"`
			Fields     map[string]any `json:"fields"`
			Object     Object         `json:"object"`
			Origin     Origin         `json:"origin"`
		}
		body := Body{
			CallbackID: callbackId,
			Fields:     testCase.InputFields,
			Object:     object,
			Origin:     origin,
		}
		serializedBody, err := json.Marshal(body)
		if err != nil {
			return nil, "", err
		}

		method := "POST"
		req, err = http.NewRequestWithContext(
			ctx,
			method,
			url,
			bytes.NewReader(serializedBody),
		)
		if err != nil {
			return nil, "", err
		}

		req.Header.Add(
			"X-HubSpot-Signature",
			signRequestV2(
				s.clientSecret,
				method,
				actionDef.Config.ActionURL,
				serializedBody,
			),
		)
	}

	req.Header.Add("Content-Type", "application/json")

	return req, callbackId, nil
}
