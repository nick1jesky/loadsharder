package client

import (
	"loadsharder/metric"
	"net/http"
	"sync"
	"time"

	"github.com/nick1jesky/atlimiter"
)

type ClientInterface interface {
	Start()
	Stop()
	Add(req *http.Request, cb Callback)
}

type Options struct {
	Client            *http.Client
	RequestTimeout    time.Duration
	ShardCount        int
	WorkersPerShard   int
	QueueSizePerShard int
	MaxRPS            uint64
	CapacityFactor    float64
}

func (o *Options) isZero() bool {
	if o.Client == nil ||
		o.RequestTimeout <= 0 ||
		o.ShardCount <= 0 ||
		o.WorkersPerShard <= 0 ||
		o.QueueSizePerShard <= 0 ||
		o.MaxRPS <= 0 ||
		o.CapacityFactor <= 0 {
		return true
	}
	return false
}

type Client struct {
	opts Options

	limiter *atlimiter.ATLimiter

	shards   []*shard
	selector ShardSelector

	startOnce sync.Once
	stopOnce  sync.Once
	wg        sync.WaitGroup

	metrics metric.Metrics
}

func NewClient(opts Options, metrics metric.Metrics) (*Client, error) {
	if opts.isZero() {
		return nil, ErrInvalidParam
	}

	return &Client{
		opts:      opts,
		limiter:   atlimiter.NewLimiter(opts.MaxRPS, opts.CapacityFactor),
		selector:  NewRoundRobinSelector(opts.QueueSizePerShard),
		startOnce: sync.Once{},
		stopOnce:  sync.Once{},
		metrics:   metrics,
	}, nil
}

func (c *Client) Start() {
	c.startOnce.Do(func() {
		c.shards = make([]*shard, c.opts.ShardCount)
		for i := range c.opts.ShardCount {
			c.shards[i] = newShard(i, c.opts, c.limiter)
			c.shards[i].start(&sync.WaitGroup{})
		}
	})
}

func (c *Client) Add(req *http.Request, cb Callback) {
	shardID := c.selector.SelectShard(c.shards)
	err := c.shards[shardID].submit(req, cb)
	if err != nil {
		cb(nil, err)
	}
}

func (c *Client) Stop() {
	c.stopOnce.Do(func() {
		for _, shard := range c.shards {
			shard.stop()
		}
		c.wg.Wait()
	})
}
