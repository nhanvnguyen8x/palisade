package log

import (
	"context"
	"log"

	"github.com/nhanvnguyen8x/palisade/internal/platform/event"
)

type Publisher struct{}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(ctx context.Context, e event.Event) error {
	log.Printf("event published: name=%s occurred_at=%s", e.Name(), e.OccurredAt())
	return nil
}
