package stats

type Snapshot interface{}

type Metrics interface {
	Snapshot() Snapshot
}

type EmptyMetrics struct{}

func (e *EmptyMetrics) Snapshot() Snapshot {
	return nil
}
