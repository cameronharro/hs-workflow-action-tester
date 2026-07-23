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
	type Origin struct {
		PortalID int `json:"portalId"`
	}
	type Object struct {
		ObjectID   int    `json:"objectId"`
		ObjectType string `json:"objectType"`
	}

	callbackId = fmt.Sprintf("ap-%d-%d-0-1", rand.Int(), rand.Int())

	method := "POST"
	url := actionDef.Config.ActionURL
	if testCase.ActionURL != "" {
		url = testCase.ActionURL
	}
	headers := map[string]string{}
	var serializedBody []byte

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

		method = string(jsCallback.HttpMethod)
		url = jsCallback.WebhookURL
		if jsCallback.HttpHeaders != nil {
			headers = jsCallback.HttpHeaders
		}
		headers["content-type"] = jsCallback.ContentType
		headers["accept"] = jsCallback.Accept

		serializedBody, err = json.Marshal(jsCallback.Body)
		if err != nil {
			return nil, "", err
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
			Fields:     testCase.InputFields,
			Object:     object,
			Origin:     origin,
		}
		serializedBody, err = json.Marshal(body)
		if err != nil {
			return nil, "", err
		}

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

	return req, callbackId, nil
}
