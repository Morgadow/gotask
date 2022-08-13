package gotask

// Runnable Interface for all subtasks. This interface must be implemented for all tasks used
type Runnable interface {
	Run()                  // runs the worker, once started no more subtasks can be added
	GetName() string       // returns name of task
	GetState() State       // returns state of task
	GetProgress() Progress // returns current progress of task
	GetWeight() Weight     // returns task weighting
	GetDesc() string       // returns task description
	GetWorkLoad() int      // returns task workload (progress times weight)
	Reset() error          // Resets task to start state
}
