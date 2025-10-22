package client

import (
	"loadsharder/ext"
	"net/http"
)

type shard struct {
	id      int
	queue   chan *http.Request
	workers []*worker
}

func newShard(opts Options, id int) *shard {
	panic(ext.ErrNotImplemented)
}

func (s *shard) start() {
	panic(ext.ErrNotImplemented)
}

func (s *shard) submit(req *http.Request) {
	panic(ext.ErrNotImplemented)
}
