package hsserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
	"github.com/cameronharro/hs-workflow-tester/internal/testcases"
)

type HSServer struct {
	clientSecret string
	actionURL    string
}

func NewHSServer(clientSecret, actionURL string) *HSServer {
	return &HSServer{
		clientSecret: clientSecret,
		actionURL:    actionURL,
	}
}

func (s *HSServer) SendRequest(actionDef actiondefinition.ActionDefinition, testCase testcases.TestCase) error {
	reqBody, err := createBody(actionDef, testCase)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	method := "POST"
	req, err := http.NewRequestWithContext(ctx, method, s.actionURL, reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("X-HubSpot-Signature", signRequestV2(s.clientSecret, method, s.actionURL, reqBody.String()))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(resBody))
	return nil
}

func createBody(actionDef actiondefinition.ActionDefinition, testCase testcases.TestCase) (*bytes.Buffer, error) {
	type Origin struct {
		PortalID int `json:"portalId"`
	}
	type Object struct {
		ObjectID   int    `json:"objectId"`
		ObjectType string `json:"objectType"`
	}
	type Body struct {
		CallbackID string         `json:"callbackId"`
		Origin     Origin         `json:"origin"`
		Object     Object         `json:"object"`
		Fields     map[string]any `json:"fields"`
	}

	fields := map[string]any{}
	for _, inputField := range testCase.InputFields {
		name := inputField.Name
		fields[name] = inputField.Value
	}
	body := Body{
		CallbackID: fmt.Sprintf("ap-%d-%d-0-1", rand.Int(), rand.Int()),
		Fields:     fields,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}
