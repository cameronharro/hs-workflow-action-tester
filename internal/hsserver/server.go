package hsserver

import (
	"sync/atomic"
	"time"
)

type HSServer struct {
	clientSecret    string
	resolutionQueue resolutionQueue
	ResultChan      <-chan error
	testsInitiated  *atomic.Int64
}

func NewHSServer(clientSecret string, port int) *HSServer {
	resultChan := make(chan error)
	testsInitiated := atomic.Int64{}
	resQueue := startResolutionQueue(&testsInitiated, 5*time.Second)
	callbackListener := startCallbackListener(port, resQueue.ResponseChan)

	go func() {
		result := <-resQueue.ResultChan
		callbackListener.Close()
		resultChan <- result
	}()

	return &HSServer{
		clientSecret:    clientSecret,
		resolutionQueue: resQueue,
		ResultChan:      resultChan,
		testsInitiated:  &testsInitiated,
	}
}
