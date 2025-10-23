package client

import (
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/nick1jesky/atlimiter"
)

type shard struct {
	id int

	limiter *atlimiter.ATLimiter

	opts  Options
	queue chan *request

	queueLen atomic.Int32

	workers  []*worker
	stopChan chan struct{}
	mu       sync.Mutex
}

func newShard(id int, opts Options, limiter *atlimiter.ATLimiter) *shard {
	return &shard{
		id:      id,
		opts:    opts,
		limiter: limiter,

		queue:    make(chan *request, opts.QueueSizePerShard),
		workers:  make([]*worker, 0, opts.WorkersPerShard),
		stopChan: make(chan struct{}),
	}
}

func (s *shard) addWorkers(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for range count {
		worker := newWorker(s.opts.Client, s.limiter)
		s.workers = append(s.workers, worker)
	}
}

// func (s *shard) removeWorkers(count int) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if count > len(s.workers) {
// 		count = len(s.workers)
// 	}
// 	for i := 0; i < count; i++ {
// 		idx := len(s.workers) - 1 - i
// 		s.workers[idx].stop()
// 	}

// 	s.workers = s.workers[:len(s.workers)-count]
// }

// func (s *shard) changeQueueLen(newSize int) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	newQueue := make(chan *request, newSize)

// 	go func() {
// 		defer close(newQueue)
// 		for req := range s.queue {
// 			newQueue <- req
// 		}
// 	}()

// 	s.queue = newQueue
// }

// func (s *shard) changeWorkersCount(newCount int) {
// 	currentCount := len(s.workers)
// 	delta := newCount - currentCount

// 	if delta > 0 {
// 		s.addWorkers(delta)
// 	} else if delta < 0 {
// 		s.removeWorkers(-delta)
// 	}
// }

// func (s *shard) setOptions(newOpts Options) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	s.opts = newOpts
// 	s.changeQueueLen(newOpts.QueueSizePerShard)
// 	s.changeWorkersCount(newOpts.WorkersPerShard)
// }

func (s *shard) start(wg *sync.WaitGroup) {
	s.addWorkers(s.opts.WorkersPerShard)

	for _, wr := range s.workers {
		wg.Add(1)
		go func(wr *worker) {
			defer wg.Done()
			wr.process(s.queue, s.stopChan)
		}(wr)
	}
}

func (s *shard) submit(req *http.Request, cb Callback) error {
	select {
	case s.queue <- &request{req: req, cb: cb}:
		s.queueLen.Add(1)
		return nil
	case <-s.stopChan:
		return ErrClientStopped
	default:
		return ErrQueueFull
	}
}

func (s *shard) stop() {
	close(s.stopChan)
}
