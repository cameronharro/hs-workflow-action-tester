package hsserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func startCallbackListener(s *HSServer, port int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /automation/actions/callbacks/2026-03/{callbackId}/complete", singleCallbackHandleFunc(s.resolutionQueue.responseChan))
	mux.HandleFunc("POST /automation/actions/callbacks/2026-03/complete", batchCallbackHandleFunc(s.resolutionQueue.responseChan))
	server := http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", port),
	}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err.Error())
			s.close()
		}
	}()
	return &server
}

func singleCallbackHandleFunc(responseChan chan<- testResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		callbackId := r.PathValue("callbackId")
		if callbackId == "" {
			w.WriteHeader(400)
			return
		}

		callbackBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		var callbackBody map[string]any
		if err := json.Unmarshal(callbackBytes, &callbackBody); err != nil {
			w.WriteHeader(400)
			return
		}

		callback := testResponse{
			CallbackId:   callbackId,
			ResponseBody: callbackBody,
		}
		responseChan <- callback
	}
}

func batchCallbackHandleFunc(responseChan chan<- testResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		callbacksBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		type CallbacksBody struct {
			Inputs []map[string]any `json:"inputs"`
		}
		var callbacksBody CallbacksBody
		if err := json.Unmarshal(callbacksBytes, &callbacksBody); err != nil {
			w.WriteHeader(400)
			return
		}

		for _, callback := range callbacksBody.Inputs {
			callbackId, ok := callback["callbackId"]
			if !ok {
				continue
			}
			callbackIdStr, ok := callbackId.(string)
			if !ok {
				continue
			}
			testRes := testResponse{
				CallbackId:   callbackIdStr,
				ResponseBody: callback,
			}
			responseChan <- testRes
		}
	}
}
