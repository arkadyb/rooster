package rooster

import (
	"errors"

	"github.com/gofrs/uuid"
)

// QueueSelector decribes strategy used to select Queue
type QueueSelector func([]*Queue) *Queue

// Mux type stores data of available queues and manages queue selection to schedule job
type Mux struct {
	queues   []*Queue
	jobsmap  map[uuid.UUID]*Queue
	selector func([]*Queue) *Queue
}

// NewMux returns new Mux item
func NewMux(selector QueueSelector, queues []*Queue) *Mux {
	return &Mux{
		jobsmap:  make(map[uuid.UUID]*Queue),
		selector: selector,
		queues:   queues,
	}
}

// Dequeue removes job from being scheduled
func (m *Mux) Dequeue(job Job) error {
	if _, ok := m.jobsmap[job.GetID()]; ok {
		queue := m.jobsmap[job.GetID()]
		return queue.Dequeue(job)
	}

	return errors.New("job cant be found")
}

// Enqueue schedules passed job
func (m *Mux) Enqueue(job *Job) {
	q := m.selector(m.queues)
	q.Enqueue(job)
	m.jobsmap[job.GetID()] = q
}
