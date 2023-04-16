package gotask

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrWorkerRunning        error = errors.New("worker already running, can not change taskqueue")
	ErrWorkerNotStarted     error = errors.New("worker was not started and is still in waiting state")
	ErrWorkerNotRunning     error = errors.New("worker is not running")
	ErrWorkerTaskQueueEmpty error = errors.New("worker task queue is empty")
	ErrWokerFinished        error = errors.New("worker already finished")
	ErrWorkerTimeoutReached error = errors.New("worker reached timeout limit")
	ErrWorkerCanceledByUser error = errors.New("worker was canceled by user")
)

// Worker Main handler struct containing all tasks and handling their run with progress evaluation
type Worker struct {
	name           string
	state          State
	progress       Progress
	taskQueue      []Runnable
	currSubTask    Runnable
	currSubTaskIdx int
	startTime      time.Time // time the worker was started
	timeoutTime    time.Time // time the timeout will be reached, if no timeout set, this is not set
	timeoutSet     bool
	wg             sync.WaitGroup // waitgroup so main routine can wait until worker is finished, used by Wait()
	quit           chan (bool)    // channel to hold quit information, this will stop running next task in line
	err            error          // return error for wait method
}

// NewWorker Factory method for creating a new worker for proper initialition
func NewWorker(name string) *Worker {
	worker := Worker{
		name:     name,
		state:    Waiting,
		progress: MinProgress,
		wg:       sync.WaitGroup{},
	}
	return &worker
}

// AddTask Starts running all tasks
// timeout Timeout in seconds which will stop worker if reached. If not set greater zero, no timeout is set.
func (w *Worker) Run(timeout time.Duration) error {
	if w.state == Running {
		return ErrWorkerRunning
	}
	if w.state == Finished || w.state == Canceled {
		return ErrWokerFinished
	}
	if w.GetAmountSubtasks() == 0 {
		return nil
	}

	// runtime and deadline evaluation
	w.err = nil
	w.state = Running
	w.startTime = time.Now()
	if timeout > 0 {
		w.timeoutSet = true
		w.timeoutTime = w.startTime.Add(timeout)
	} else {
		w.timeoutSet = false
	}
	w.quit = make(chan bool)

	// create channel to store state in and
	w.wg.Add(1)
	go w.runInternal()

	return nil
}

// Wait Wait until worker is finished
func (w *Worker) Wait() error {
	if w.state != Running {
		return ErrWorkerNotRunning
	}
	w.wg.Wait()
	return w.err
}

// Stop Stops task run
func (w *Worker) Stop() error {
	if w.state != Running {
		return ErrWorkerNotRunning
	}
	w.quit <- true
	w.err = ErrWorkerCanceledByUser
	return nil
}

// Reset Can be used to reset worker to status quo state to run again after run once
func (w *Worker) Reset() error {
	if w.state == Running {
		return ErrWorkerRunning
	}
	w.state = Waiting
	w.progress = MinProgress
	w.err = nil
	for _, task := range w.taskQueue {
		task.Reset()
	}
	return nil
}

// AddTask Adds new task to queue
func (w *Worker) AddTask(task Runnable) error {
	if w.state == Running {
		return ErrWorkerRunning
	}
	w.taskQueue = append(w.taskQueue, task)
	return nil
}

// AddTask Adds multiple new tasks to queue
func (w *Worker) AddTasks(tasks []Runnable) error {
	var err error = nil
	for _, task := range tasks {
		err = w.AddTask(task)
		if err != nil {
			return err
		}
	}
	return err
}

// AddTask Emptys task queue
func (w *Worker) ClearTasks() error {
	if w.state == Running {
		return ErrWorkerRunning
	}
	w.taskQueue = nil
	return nil
}

// GetAmountSubtasks Returns amount of tasks in queue
func (w *Worker) GetAmountSubtasks() int {
	return len(w.taskQueue)
}

// GetName Returns present worker name
func (w *Worker) GetName() string {
	return w.name
}

// GetState Returns present worker state
func (w *Worker) GetState() State {
	return w.state
}

// GetProgress Returns present queue progress in percent from 0 to 100
func (w *Worker) GetProgress() Progress {
	if w.state == Running {
		w.updateProgress()
	}
	return w.progress
}

// GetTotalWorkLoad Returns total workload of all tasks in queue combined (progress times weight)
func (w *Worker) GetTotalWorkLoad() float64 {
	totalLoad := 0.0
	for _, task := range w.taskQueue {
		totalLoad += float64(task.GetWeight())
	}
	return totalLoad
}

// GetRemainingWorkLoad Returns remaining workload of all tasks in queue combined (progress times weight)
func (w *Worker) GetRemainingWorkLoad() float64 {
	remainLoad := 0.0
	for _, task := range w.taskQueue {
		remainLoad += (1 - float64(task.GetProgress())/float64(MaxProgress)) * float64(task.GetWeight())
	}
	return remainLoad
}

// GetDuration Get duration for how long worker was or is running in seconds
// Note: Not to be called on worker in Waiting state
func (w *Worker) GetDuration() (float64, error) {
	if w.state == Waiting {
		return 0, ErrWorkerNotStarted
	}
	return float64(time.Since(w.startTime)/time.Millisecond) / 1000, nil
}

// GetDuration Get duration for how long worker was or is running in seconds
// Note: Only to be called during running worker. If no timeout set, a -1 is returned
func (w *Worker) GetRemainingTime() (float64, error) {
	if w.state != Running {
		return 0, ErrWorkerNotRunning
	}
	if !w.timeoutSet {
		return -1, nil
	}
	return float64(time.Until(w.timeoutTime)/time.Millisecond) / 1000, nil
}

// IsReady ReConvienince function to check if worker is ready to start
func (w *Worker) IsReady() bool {
	return w.state == Waiting
}

// IsRunning ReConvienince function to check if worker currently running
func (w *Worker) IsRunning() bool {
	return w.state == Running
}

// IsFinished ReConvienince function to check if worker finished its run
func (w *Worker) IsFinished() bool {
	return w.state == Finished
}

// GetSubtasks Returns copy of all subtasks as slice
func (w *Worker) GetSubtasks() []Runnable {
	return w.taskQueue
}

// GetCurrentTaskName Returns name of presently running task
func (w *Worker) GetCurrentTaskName() (string, error) {
	if w.state != Running || w.currSubTask == nil {
		return "", ErrWorkerNotRunning
	}
	if w.GetAmountSubtasks() == 0 {
		return "", ErrWorkerTaskQueueEmpty
	}
	return w.currSubTask.GetName(), nil
}

// GetCurrentTaskDesc Returns description of presently running task
func (w *Worker) GetCurrentTaskDesc() (string, error) {
	if w.state != Running || w.currSubTask == nil {
		return "", ErrWorkerNotRunning
	}
	if w.GetAmountSubtasks() == 0 {
		return "", ErrWorkerTaskQueueEmpty
	}
	return w.currSubTask.GetDesc(), nil
}

// updateProgress Updates internal progress over all tasks
func (w *Worker) updateProgress() {
	workTotal := 0
	workDone := 0
	for _, task := range w.taskQueue {
		workTotal += int(task.GetWeight())
		workDone += int(task.GetProgress()/100) * int(task.GetWeight())
	}
	w.progress = (Progress(workDone) / Progress(workTotal)) * 100 // multiply by 100 for percent
}

// runInternal Internal run function which is run in another context to handle timeout and termination
func (w *Worker) runInternal() {
	w.state = Running
out:
	for idx, task := range w.taskQueue {
		select {
		case <-w.quit:
			w.state = Canceled
			w.err = ErrWorkerCanceledByUser
			break out
		default:
			w.updateProgress()

			if w.timeoutReached() {
				w.state = TimeoutReached
				w.err = ErrWorkerTimeoutReached
				break out
			}

			// call next subtask
			w.currSubTaskIdx = idx
			w.currSubTask = task
			w.currSubTask.Run()

			if w.currSubTaskIdx == len(w.taskQueue)-1 {
				w.state = Finished
				w.updateProgress()
			}
		}
	}
	w.wg.Done()
}

// timeoutReached Checks if timeout of worker was reached
func (w *Worker) timeoutReached() bool {
	remain, _ := w.GetRemainingTime()
	return w.timeoutSet && remain < 0
}
