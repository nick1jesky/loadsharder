package client

import "net/http"

type AsyncResult struct {
	Resp *http.Response
	Err  error
}

type request struct {
	req        *http.Request
	resultChan chan<- AsyncResult
}
