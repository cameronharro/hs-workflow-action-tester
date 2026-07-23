package hsserver

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

type HSServer struct {
	callbackListener *http.Server
	clientSecret     string
	ctx              context.Context
	cancelFunc       context.CancelFunc
	resolutionQueue  *resolutionQueue
	result           *Result
	waitChan         chan error
}

type Result struct {
	error error
	mu    sync.Mutex
}

func (s *HSServer) AddResult(error error) {
	s.result.mu.Lock()
	defer s.result.mu.Unlock()
	s.result.error = errors.Join(s.result.error, error)
	s.resolutionQueue.responsesProcessed.Add(1)
}

func NewHSServer(clientSecret string, port int, timeout time.Duration) *HSServer {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)

	server := &HSServer{
		clientSecret: clientSecret,
		ctx:          ctx,
		cancelFunc:   cancelFunc,
		result:       &Result{},
		waitChan:     make(chan error),
	}
	server.resolutionQueue = newResolutionQueue(server)
	server.callbackListener = startCallbackListener(server, port)

	go func() {
		defer cancelFunc()
		select {
		case <-ctx.Done():
			server.callbackListener.Close()
			server.waitChan <- server.result.error
		}

	}()

	return server
}

func (s *HSServer) close() {
	s.cancelFunc()
}

func (s *HSServer) Wait() error {
	return <-s.waitChan
}
