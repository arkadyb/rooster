package rooster

import (
	"errors"

	"github.com/gofrs/uuid"
)

// QueueSelector decribes strategy used to select Queue
type QueueSelectorFunc func([]*Queue) *Queue

// Rooster type stores data of available queues and manages queue selection to schedule job
type Rooster struct {
	queues   []*Queue
	jobsmap  map[uuid.UUID]*Queue
	selector func([]*Queue) *Queue
}

// NewMux returns new Rooster item
func NewRooster(selector QueueSelectorFunc, queues []*Queue) *Rooster {
	return &Rooster{
		jobsmap:  make(map[uuid.UUID]*Queue),
		selector: selector,
		queues:   queues,
	}
}

// Dequeue removes job from being scheduled
func (m *Rooster) Dequeue(job Job) error {
	if _, ok := m.jobsmap[job.GetID()]; ok {
		queue := m.jobsmap[job.GetID()]
		return queue.Dequeue(job)
	}

	return errors.New("job cant be found")
}

// Enqueue schedules passed job
func (m *Rooster) Enqueue(job *Job) {
	q := m.selector(m.queues)
	q.Enqueue(job)
	m.jobsmap[job.GetID()] = q
}
