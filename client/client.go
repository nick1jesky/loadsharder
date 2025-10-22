package client

import (
	"loadsharder/ext"
	"loadsharder/metrics"
	"net/http"
	"sync"
	"time"

	"github.com/nick1jesky/atlimiter"
)

type Options struct {
	Client            *http.Client
	Limiter           *atlimiter.ATLimiter
	RequestTimeout    time.Duration
	ShardCount        int
	WorkersPerShard   int
	QueueSizePerShard int
}

func (o *Options) IsZero() bool {
	if o.Client == nil ||
		o.Limiter == nil ||
		o.RequestTimeout == 0 ||
		o.ShardCount <= 0 ||
		o.WorkersPerShard <= 0 ||
		o.QueueSizePerShard <= 0 {
		return true
	}
	return false
}

type Client struct {
	opts Options

	shards   []*shard
	selector ShardSelector

	resChan chan *AsyncResult

	stopOnce sync.Once
	stopChan chan struct{}
	metrics  metrics.Metrics
}

func NewClient(opts Options) (*Client, error) {
	if opts.IsZero() {
		return nil, ErrInvalidParam
	}

	client := &Client{
		opts:     opts,
		stopChan: make(chan struct{}),
		selector: NewRoundRobinSelector(),
	}

	return client, nil
}

func (c *Client) Start() error {
	c.shards = make([]*shard, c.opts.ShardCount)
	for i := range c.opts.ShardCount {
		c.shards[i] = newShard(c.opts, i)
		c.shards[i].start()
	}
	return ext.ErrNotImplemented
}

func (c *Client) Add(req *http.Request) {
	if !c.opts.Limiter.Allow() {
		panic(ext.ErrNotImplemented)
	}

	shardID := c.selector.SelectShard()
	c.shards[shardID].submit(req)

	panic(ext.ErrNotImplemented)
}
