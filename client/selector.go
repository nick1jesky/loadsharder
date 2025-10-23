package client

type ShardSelector interface {
	SelectShard(sh []*shard) int
}

type RoundRobinSelector struct {
	maxQueueLen int
}

func NewRoundRobinSelector(maxQueueLen int) *RoundRobinSelector {
	return &RoundRobinSelector{
		maxQueueLen: maxQueueLen,
	}
}

func (s *RoundRobinSelector) SelectShard(sh []*shard) int {
	if sh == nil {
		return 0
	}

	min := s.maxQueueLen
	minIdx := 0

	for itt, j := range sh {
		tmp := int(j.queueLen.Load())
		if tmp <= min {
			min = tmp
			minIdx = itt
		}
	}

	return minIdx
}
