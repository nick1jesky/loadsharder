package client

import "sync/atomic"

type ShardSelector interface {
	SelectShard() int
}

type RoundRobinSelector struct {
	counter uint64
	shards  int
}

func NewRoundRobinSelector() *RoundRobinSelector {
	return &RoundRobinSelector{}
}

func (s *RoundRobinSelector) SetShardCount(count int) {
	s.shards = count
}

func (s *RoundRobinSelector) SelectShard() int {
	if s.shards == 0 {
		return 0
	}
	return int(atomic.AddUint64(&s.counter, 1) % uint64(s.shards))
}
