package client

import (
	"net/http"
)

type worker struct {
	client *http.Client
}

func newWorker(client *http.Client) *worker {
	return &worker{
		client: client,
	}
}

func (w *worker) process(queue <-chan *request) {
	for req := range queue {
		resp, err := w.client.Do(req.req)
		req.resultChan <- AsyncResult{
			Resp: resp,
			Err:  err,
		}
		close(req.resultChan)
	}
}
