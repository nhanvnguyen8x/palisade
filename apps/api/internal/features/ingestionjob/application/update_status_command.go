package application

import (
	"github.com/google/uuid"
	"github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/domain"
)

type UpdateStatusCommand struct {
	JobID        uuid.UUID
	Status       domain.JobStatus
	ErrorMessage string
}
