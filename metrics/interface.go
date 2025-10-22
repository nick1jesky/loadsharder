package metrics

type Metrics interface {
	Snapshot() Snapshot
}
