package hsserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/jshelper"
	"github.com/cameronharro/hs-workflow-tester/internal/testcase"
)

func (s *HSServer) CreateRequest(
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
	type Body struct {
		CallbackID string         `json:"callbackId"`
		Fields     map[string]any `json:"fields"`
		Object     Object         `json:"object"`
		Origin     Origin         `json:"origin"`
	}

	callbackId = fmt.Sprintf("ap-%d-%d-0-1", rand.Int(), rand.Int())
	body := Body{
		CallbackID: callbackId,
		Fields:     testCase.InputFields,
		Object: Object{
			ObjectID:   testCase.ObjectID,
			ObjectType: testCase.ObjectType,
		},
		Origin: Origin{
			PortalID: testCase.PortalID,
		},
	}
	serializedBody, err := json.Marshal(body)
	if err != nil {
		return nil, "", err
	}

	if preActionFunction := actionDef.GetPreActionFunction(); preActionFunction != nil {
		jsCallback, err := jshelper.RunPreActionFunction(
			jshelper.PreActionEvent{
				InputFields: testCase.InputFields,
				WebhookURL:  actionDef.Config.ActionURL,
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

		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
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
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		method := "POST"
		req, err = http.NewRequestWithContext(
			ctx,
			method,
			actionDef.Config.ActionURL,
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
