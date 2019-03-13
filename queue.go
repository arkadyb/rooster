package rooster

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var unixInifiniteDuration = time.Until(time.Unix(1<<63-62135596801, 999999999))

// QueueInterceptor describes middleware function run for the job
type QueueInterceptor func(job Job)

// Queue types contains all data of scheduled jobs to be run in given queue
type Queue struct {
	lock sync.Mutex
	jobs []*Job

	interceptors []QueueInterceptor
	timer        *time.Timer
}

// NewQueue returns new Queue item
func NewQueue(interceptors ...QueueInterceptor) *Queue {
	q := &Queue{
		jobs: make([]*Job, 0),

		timer:        time.NewTimer(unixInifiniteDuration),
		interceptors: interceptors,
	}

	go func(q *Queue) {
		for {
			<-q.timer.C

			q.lock.Lock()
			j := q.jobs[0]

			if len(q.jobs) > 1 {
				// remove 0th job frm the slice
				q.jobs = q.jobs[1:]
				q.timer.Reset(time.Until(q.jobs[0].when))
			} else {
				// turn queue into idle mode
				q.timer.Reset(unixInifiniteDuration)
				q.jobs = []*Job{}
			}
			q.lock.Unlock()

			// run job
			if len(q.interceptors) > 0 {
				var wg sync.WaitGroup
				for _, interceptor := range q.interceptors {
					wg.Add(1)
					go func(i QueueInterceptor) {
						i(*j)
						wg.Done()
					}(interceptor)
				}
				wg.Wait()
			}
			go j.f(j)
		}

	}(q)

	return q
}

// Dequeue removes job from the queue
func (q *Queue) Dequeue(job Job) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	for i, j := range q.jobs {
		if j.GetID() == job.GetID() {

			if len(q.jobs) > 1 {
				q.jobs = append(q.jobs[:i], q.jobs[i+1:]...)
				q.timer.Reset(time.Until(q.jobs[0].when))
			} else {
				q.jobs = q.jobs[:0]
				q.timer.Reset(unixInifiniteDuration)
			}

			return nil
		}
	}

	return errors.New("job cant be found")
}

// Enqueue inserts queue into scheduled job list
func (q *Queue) Enqueue(job *Job) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.jobs = insert(q.jobs, job)
	if q.jobs[0].GetID() == job.GetID() {
		q.timer.Reset(time.Until(job.when))
	}
}

func insert(data []*Job, el *Job) []*Job {
	index := sort.Search(len(data), func(i int) bool { return data[i].when.After(el.when) })
	data = append(data, &Job{})
	copy(data[index+1:], data[index:])
	data[index] = el
	return data
}
