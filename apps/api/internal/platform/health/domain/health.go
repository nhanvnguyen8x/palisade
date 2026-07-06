package domain

type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

type HealthReport struct {
	Status Status
}
