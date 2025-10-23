package client

type worker struct {
	shardPtr *shard
	stopChan chan struct{}
}

func newWorker(sh *shard) *worker {
	return &worker{
		shardPtr: sh,
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
		if w.shardPtr.limiter.Allow() {
			resp, err := w.shardPtr.opts.Client.Do(req.req)
			req.cb(resp, err)
			w.shardPtr.queueLen.Add(-1)
			return
		}
	}
}
