package dispatcher

import (
	"fmt"
	"log"
)

type Dispatcher struct {
	TaskCounter    int                  // internal counter for number of tasks
	taskQueue      chan *Task           // channel of tasks submitted by main()
	dispatchStatus chan *DispatchStatus // channel for task/worker status reports
	workQueue      chan *Task           // channel of work dispatched
	workerQueue    chan *Worker         // channel of workers
}

type Worker struct {
	ID             int
	tasks          chan *Task
	dispatchStatus chan *DispatchStatus
	Quit           chan bool
}

type Task struct {
	ID int
	F  func() error
}

type DispatchStatus struct {
	Type   string
	ID     int
	Status string
}

type TaskExecutable func() error

func CreateNewDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{
		TaskCounter:    0,
		taskQueue:      make(chan *Task),
		dispatchStatus: make(chan *DispatchStatus),
		workQueue:      make(chan *Task),
		workerQueue:    make(chan *Worker),
	}
	return dispatcher
}

func (dispatcher *Dispatcher) AddTask(te TaskExecutable) {
	task := &Task{
		ID: dispatcher.TaskCounter,
		F:  te,
	}
	go func() {
		log.Println("GIRDI")
		log.Println(dispatcher.taskQueue)
		dispatcher.taskQueue <- task
		log.Println("CIKTI")
		log.Println(dispatcher.taskQueue)
	}()
	dispatcher.TaskCounter++
	fmt.Printf("TaskCounter is now: %d\n", dispatcher.TaskCounter)
}

func CreateNewWorker(id int, workerQueue chan *Worker, taskQueue chan *Task, dStatus chan *DispatchStatus) *Worker {
	w := &Worker{
		ID:             id,
		tasks:          taskQueue,
		dispatchStatus: dStatus,
	}
	go func() {
		workerQueue <- w
	}()
	return w
}

func (dispatcher *Dispatcher) Start(numWorkers int) {
	// Create numworkers
	for i := 0; i < numWorkers; i++ {
		worker := CreateNewWorker(i, dispatcher.workerQueue, dispatcher.workQueue, dispatcher.dispatchStatus)
		worker.Start()
	}

	// wait for work to be added then pass it off.
	go func() {
		for {
			select {
			case task := <-dispatcher.taskQueue:
				log.Printf("Got a task in the queue to dispatch: %d\n", task.ID)
				// Sending it off;
				dispatcher.workQueue <- task
			case dStatus := <-dispatcher.dispatchStatus:
				log.Printf("Got a dispatch status: \n\tType[%s] - ID[%d] - Status[%s]\n", dStatus.Type, dStatus.ID, dStatus.Status)
				if dStatus.Type == "worker" {
					if dStatus.Status == "quit" {
						dispatcher.TaskCounter--
					}
				}
			}

		}
	}()

}
func (worker *Worker) Start() {
	go func() {
		for {
			select {
			case task := <-worker.tasks:
				log.Printf("Worker[%d] executing task[%d]\n.", worker.ID, task.ID)
				task.F()
				worker.dispatchStatus <- &DispatchStatus{
					Type:   "worker",
					ID:     worker.ID,
					Status: "quit",
				}
				worker.Quit <- true
			case <-worker.Quit:
				return
			}
		}
	}()
}

func (dispatcher *Dispatcher) Finished() bool {
	if dispatcher.TaskCounter < 1 {
		return true
	}
	return false
}
