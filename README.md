# rooster

Rooster is a job scheduler that runs in multiple queues with jobs being executed in parallel.

# How to use

```go
import "github.com/arkadyb/rooster"
```

## Create new rooster

To create new rooster item use `NewRooster` function that takes two arguments - `selector` and `queues` where `selector` is function of type `QueueSelectorFunc` that describes a way next execution queue is selected for a new scheduled job.

You may define your own `QueueSelectorFunc` to select queue depending on content of your job or scheduler.
A simple Round Robin selector may look like  
```go
roundRobinSelector:=func() rooster.QueueSelectorFunc {
    selectedQueueIndex := 0
    return func(queues []*rooster.Queue) *rooster.Queue {
        q := queues[selectedQueueIndex]

        selectedQueueIndex++
        if selectedQueueIndex >= len(queues) {
            selectedQueueIndex = 0
        }

        return q
    }
}
```

`Queues` are the workload chains that store and execute jobs. User `rooster.NewQueue()` function to create new queue.

```go
scheduler:=rooster.NewRooster(roundRobinSelector, []*rooster.Queue{rooster.NewQueue(), rooster.NewQueue()})
```
Here rooster would be created with two execution queues with one being selected to run scheduled job in round robin style.  

## Schedule job

Rooster exposes `Enqueue` function to schedule new job and `Dequeue` to remove job before its being executed. Use `rooster.NewJob()` function to create a new job. 
Function takes two arguments - `time` when job should be executed and `JobFunc`. 

```go
scheduler.Enqueue(NewJob(time.Now().Add(5 * time.Second), func(j *Job){
    fmt.Printf("Job %s has been executed\n", j.GetID())
}))
```   

## Middleware 

Interceptors may be included into execution chain for given `Queue`. Specify interceptor functions when creating new queue
```go
NewQueue(func(job Job){
    fmt.Printf("First interceptor run for job %s\n", job.GetID())
}, func(job Job){
    fmt.Printf("Second interceptor run for job %s\n", job.GetID())
})
```

## License
 
The MIT License (MIT)
