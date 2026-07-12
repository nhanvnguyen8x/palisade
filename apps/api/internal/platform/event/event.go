package event

import "time"

type Event interface {
	Name() string
	OccurredAt() time.Time
}
