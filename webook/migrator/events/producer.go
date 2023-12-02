package events

import "context"

//go:generate mockgen -source=producer.go -package=evtmocks -destination=mocks/producer.mock.go Producer
type Producer interface {
	ProduceInconsistentEvent(ctx context.Context, event InconsistentEvent) error
}
