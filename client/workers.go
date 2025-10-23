package client

import (
	"net/http"
	"time"

	"github.com/nick1jesky/atlimiter"
)

type worker struct {
	client   *http.Client
	limiter  *atlimiter.ATLimiter
	stopChan chan struct{}
}

func newWorker(client *http.Client, limiter *atlimiter.ATLimiter) *worker {
	return &worker{
		client:   client,
		limiter:  limiter,
		stopChan: make(chan struct{}),
	}
}

func (w *worker) process(queue <-chan *request, shardStopChan <-chan struct{}) {
	for {
		select {
		case req, ok := <-queue:
			if !ok {
				return
			}
			w.processRequest(req)
		case <-shardStopChan:
			return
		case <-w.stopChan:
			return
		}
	}
}

func (w *worker) processRequest(req *request) {
	for {
		switch {
		case w.limiter.Allow():
			resp, err := w.client.Do(req.req)
			req.cb(resp, err)
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
