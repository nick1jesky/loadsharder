package client

import "net/http"

type Callback func(*http.Response, error)

type request struct {
	req *http.Request
	cb  Callback
}
