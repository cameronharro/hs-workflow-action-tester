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
	return &HSServer{
		clientSecret:    clientSecret,
		resolutionQueue: startResolutionQueue(&testsInitiated, resultChan, 5*time.Second),
		ResultChan:      resultChan,
		testsInitiated:  &testsInitiated,
	}
}
