package rooster

import (
	"time"

	"github.com/gofrs/uuid"
)

// JobFunc used to describe functionality to to run by the scheduler
type JobFunc func(*Job)

// Job describes work item used in Queues and Mux
type Job struct {
	id uuid.UUID

	when time.Time
	f    JobFunc
}

// NewJob returns new Job item
func NewJob(when time.Time, f JobFunc) *Job {
	return &Job{
		id: uuid.NewV4(),

		when: when,
		f:    f,
	}
}

// GetID returns given job ID
func (s Job) GetID() uuid.UUID {
	return s.id
}

// Func returns given job work function
func (s Job) Func() JobFunc {
	return s.f
}
