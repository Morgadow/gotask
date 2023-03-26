package test

import (
	"testing"
	"time"

	"github.com/morgadow/gotask"
)

// Sleeping test function which halts goroutine for defined length
func Sleeping(dur interface{}) {
	duration := dur.(int)
	time.Sleep(time.Millisecond * time.Duration(duration))
}

// helper function
func createWorker() *gotask.Worker {
	worker := gotask.NewWorker("Workername")
	_ = worker.AddTask(gotask.NewTask("task 0", gotask.Weight(1), "Sleeping for 50ms", Sleeping, 50))
	_ = worker.AddTask(gotask.NewTask("task 1", gotask.Weight(2), "Sleeping for 50ms", Sleeping, 50))
	_ = worker.AddTask(gotask.NewTask("task 2", gotask.Weight(3), "Sleeping for 50ms", Sleeping, 50))

	return worker
}

func TestWorkerGetName(t *testing.T) {

	worker := gotask.NewWorker("Workername")
	if name := worker.GetName(); name != "Workername" {
		t.Errorf("Name of worker not equal to 'Workername': %v", name)
	}
}

func TestAddTask(t *testing.T) {

	worker := gotask.NewWorker("Workername")

	err := worker.AddTask(gotask.NewTask("task 0", gotask.Weight(1), "Sleeping for 25ms", Sleeping, 25))
	if err != nil {
		t.Errorf("err not nil: %v", err)
	}
	err = worker.AddTask(gotask.NewTask("task 1", gotask.Weight(2), "Sleeping for 50ms", Sleeping, 50))
	if err != nil {
		t.Errorf("err not nil: %v", err)
	}
	err = worker.AddTask(gotask.NewTask("task 2", gotask.Weight(3), "Sleeping for 75ms", Sleeping, 75))
	if err != nil {
		t.Errorf("err not nil: %v", err)
	}

	if amountTask := worker.GetAmountSubtasks(); amountTask != 3 {
		t.Errorf("Amount of subtasks not equal to 3: %v", amountTask)
	}
}

func TestGetAmountSubtasks(t *testing.T) {

	worker := createWorker()
	if amountTask := worker.GetAmountSubtasks(); amountTask != 3 {
		t.Errorf("Amount of subtasks not equal to 3: %v", amountTask)
	}
}

func TestAddTasks(t *testing.T) {

	worker := gotask.NewWorker("Workername")

	tasks := []gotask.Runnable{
		gotask.NewTask("task 0", gotask.Weight(1), "Sleeping for 25ms", Sleeping, 25),
		gotask.NewTask("task 1", gotask.Weight(2), "Sleeping for 50ms", Sleeping, 50),
		gotask.NewTask("task 2", gotask.Weight(3), "Sleeping for 75ms", Sleeping, 75),
	}

	err := worker.AddTasks(tasks)
	if err != nil {
		t.Errorf("err not nil: %v", err)
	}
	if amountTask := worker.GetAmountSubtasks(); amountTask != 3 {
		t.Errorf("Amount of subtasks not equal to 3: %v", amountTask)
	}
}

func TestClearTasks(t *testing.T) {

	worker := createWorker()
	if amountTask := worker.GetAmountSubtasks(); amountTask != 3 {
		t.Errorf("Amount of subtasks not equal to 3: %v", amountTask)
	}

	worker.ClearTasks()
	if amountTask := worker.GetAmountSubtasks(); amountTask != 0 {
		t.Errorf("Amount of subtasks not equal to 0: %v", amountTask)
	}
}

func TestGetTotalWorkLoad(t *testing.T) {

	worker := gotask.NewWorker("Workername")

	_ = worker.AddTask(gotask.NewTask("task 0", gotask.Weight(1), "Sleeping for 25ms", Sleeping, 25))
	_ = worker.AddTask(gotask.NewTask("task 1", gotask.Weight(2), "Sleeping for 50ms", Sleeping, 50))
	_ = worker.AddTask(gotask.NewTask("task 2", gotask.Weight(3), "Sleeping for 75ms", Sleeping, 75))
	_ = worker.AddTask(gotask.NewTask("task 3", gotask.Weight(4), "Sleeping for 100ms", Sleeping, 100))
	if weight := worker.GetTotalWorkLoad(); weight != 10 {
		t.Errorf("total weight not equal to 10: %v", weight)
	}

	worker.Run(0)
	worker.Wait()
	if weight := worker.GetTotalWorkLoad(); weight != 10 {
		t.Errorf("total weight not equal to 10: %v", weight)
	}
}

func TestGetRemainingWorkLoad(t *testing.T) {

	worker := gotask.NewWorker("Workername")

	_ = worker.AddTask(gotask.NewTask("task 0", gotask.Weight(1), "Sleeping for 25ms", Sleeping, 25))
	_ = worker.AddTask(gotask.NewTask("task 1", gotask.Weight(2), "Sleeping for 50ms", Sleeping, 50))
	_ = worker.AddTask(gotask.NewTask("task 2", gotask.Weight(3), "Sleeping for 75ms", Sleeping, 75))
	_ = worker.AddTask(gotask.NewTask("task 3", gotask.Weight(4), "Sleeping for 100ms", Sleeping, 100))
	if weight := worker.GetRemainingWorkLoad(); weight != 10 {
		t.Errorf("remaining weight not equal to 10: %v", weight)
	}

	worker.Run(0)
	worker.Wait()
	if weight := worker.GetRemainingWorkLoad(); weight != 0 {
		t.Errorf("remaining weight not equal to 0: %v", weight)
	}
}

func TestRun(t *testing.T) {

	worker := createWorker()
	if state := worker.GetState(); state != gotask.Waiting {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Waiting), gotask.StateToString(state))
	}

	err := worker.Run(0)
	if err != nil {
		t.Errorf("err not nil: %v", err)
	}
	if state := worker.GetState(); state != gotask.Running {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Running), gotask.StateToString(state))
	}

	worker.Wait()
	if state := worker.GetState(); state != gotask.Finished {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(state))
	}
}

func TestGetDuration(t *testing.T) {

	worker := createWorker()
	_ = worker.Run(0)
	worker.Wait()

	dur, err := worker.GetDuration()
	if err != nil {
		t.Errorf("err not nil: %v", err)
	}
	if dur < 0.150 || dur > 0.152 {
		t.Errorf("duration not 0.150: %v", dur)
	}
}

func TestGetState(t *testing.T) {

	worker := createWorker()
	if state := worker.GetState(); state != gotask.Waiting {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Waiting), gotask.StateToString(state))
	}

	worker.Run(0)
	if state := worker.GetState(); state != gotask.Running {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Running), gotask.StateToString(state))
	}

	worker.Wait()
	if state := worker.GetState(); state != gotask.Finished {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(state))
	}

	worker.Reset()
	worker.Run(0)
	worker.Stop()
	if state := worker.GetState(); state != gotask.Canceled {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Canceled), gotask.StateToString(state))
	}
}

func TestGetStateBool(t *testing.T) {

	worker := createWorker()
	if state := worker.GetState(); state != gotask.Waiting {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Waiting), gotask.StateToString(state))
	}
	if !worker.IsReady() {
		t.Errorf("Expected worker state IsReady to be true, got %v", worker.IsReady())
	}

	worker.Run(0)
	if state := worker.GetState(); state != gotask.Running {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Running), gotask.StateToString(state))
	}
	if !worker.IsRunning() {
		t.Errorf("Expected worker state IsRunning to be true, got %v", worker.IsRunning())
	}

	worker.Wait()
	if state := worker.GetState(); state != gotask.Finished {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(state))
	}
	if !worker.IsFinished() {
		t.Errorf("Expected worker state IsFinished to be true, got %v", worker.IsFinished())
	}

	worker.Reset()
	worker.Run(0)
	worker.Stop()
	if state := worker.GetState(); state != gotask.Canceled {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Canceled), gotask.StateToString(state))
	}
}

func TestWait(t *testing.T) {

	worker := createWorker()
	if state := worker.GetState(); state != gotask.Waiting {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Waiting), gotask.StateToString(state))
	}

	worker.Run(0)
	Sleeping(250)
	if state := worker.GetState(); state != gotask.Finished {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(state))
	}

	worker.Reset()
	if state := worker.GetState(); state != gotask.Waiting {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Waiting), gotask.StateToString(state))
	}

	worker.Run(0)
	worker.Wait()
	if state := worker.GetState(); state != gotask.Finished {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 0 {
		t.Errorf("remaining weight not equal to 0: %v", weight)
	}
}

func TestStop(t *testing.T) {

	worker := createWorker()

	// stopping immediately will set second task to not done, as the stop call is called after the first task was finished
	worker.Run(0)
	time.Sleep(1 * time.Millisecond)
	worker.Stop()
	if state := worker.GetState(); state != gotask.Canceled {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Canceled), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 5 {
		t.Errorf("remaining weight not equal to 5: %v", weight)
	}

	// stopping after first task finished will set remaining weight to 1 as first two subtasks will be finished
	worker.Reset()
	worker.Run(0)
	time.Sleep(75 * time.Millisecond)
	worker.Stop()
	if state := worker.GetState(); state != gotask.Canceled {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Canceled), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 3 {
		t.Errorf("remaining weight not equal to 3: %v", weight)
	}

	// finishing after all are done has no effect at all and returns err
	worker.Reset()
	worker.Run(0)
	time.Sleep(750 * time.Millisecond)
	worker.Stop()
	if state := worker.GetState(); state != gotask.Finished {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 0 {
		t.Errorf("remaining weight not equal to 0: %v", weight)
	}
}

func TestSubtaskState(t *testing.T) {

	worker := createWorker()

	// stopping immediately will set second task to not done, as the stop call is called after the first task was finished
	worker.Run(0)
	time.Sleep(1 * time.Millisecond)
	worker.Stop()
	if state := worker.GetState(); state != gotask.Canceled {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Canceled), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 5 {
		t.Errorf("remaining weight not equal to 5: %v", weight)
	}
	subTasks := worker.GetSubtasks()
	if subTasks[0].GetState() != gotask.Finished || subTasks[1].GetState() != gotask.Waiting || subTasks[2].GetState() != gotask.Waiting {
		t.Errorf("expected state task 1 %v, task 2 %v, task 3 %v, got: 1: %v, 2: %v, 3: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(gotask.Waiting), gotask.StateToString(gotask.Waiting), gotask.StateToString(subTasks[0].GetState()), gotask.StateToString(subTasks[1].GetState()), gotask.StateToString(subTasks[2].GetState()))
	}

	// stopping after first task finished will set remaining weight to 1 as first two subtasks will be finished
	worker.Reset()
	worker.Run(0)
	time.Sleep(75 * time.Millisecond)
	worker.Stop()
	if state := worker.GetState(); state != gotask.Canceled {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Canceled), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 3 {
		t.Errorf("remaining weight not equal to 3: %v", weight)
	}
	subTasks = worker.GetSubtasks()
	if subTasks[0].GetState() != gotask.Finished || subTasks[1].GetState() != gotask.Finished || subTasks[2].GetState() != gotask.Waiting {
		t.Errorf("expected state task 1 %v, task 2 %v, task 3 %v, got: 1: %v, 2: %v, 3: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(gotask.Finished), gotask.StateToString(gotask.Waiting), gotask.StateToString(subTasks[0].GetState()), gotask.StateToString(subTasks[1].GetState()), gotask.StateToString(subTasks[2].GetState()))
	}

	// finishing after all are done has no effect at all and returns err
	worker.Reset()
	worker.Run(0)
	time.Sleep(750 * time.Millisecond)
	worker.Stop()
	if state := worker.GetState(); state != gotask.Finished {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 0 {
		t.Errorf("remaining weight not equal to 0: %v", weight)
	}
	if subTasks[0].GetState() != gotask.Finished || subTasks[1].GetState() != gotask.Finished || subTasks[2].GetState() != gotask.Finished {
		t.Errorf("expected state task 1 %v, task 2 %v, task 3 %v, got: 1: %v, 2: %v, 3: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(gotask.Finished), gotask.StateToString(gotask.Finished), gotask.StateToString(subTasks[0].GetState()), gotask.StateToString(subTasks[1].GetState()), gotask.StateToString(subTasks[2].GetState()))
	}
}

func TestTimeout(t *testing.T) {

	worker := createWorker()
	if state := worker.GetState(); state != gotask.Waiting {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Waiting), gotask.StateToString(state))
	}

	worker.Run(75 * time.Millisecond)
	if state := worker.GetState(); state != gotask.Running {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Running), gotask.StateToString(state))
	}
	timeLeft, err := worker.GetRemainingTime()
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}
	if timeLeft > 0.0750 || timeLeft < 0.070 {
		t.Errorf("invalid timeout, expected: %v s, got: %v s", 0.075, timeLeft)
	}

	// sleep, tasks should be canceled by timeout as total worker time is 150 ms and timeout set to 75
	time.Sleep(250 * time.Millisecond)
	if state := worker.GetState(); state != gotask.TimeoutReached {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.TimeoutReached), gotask.StateToString(state))
	}

	// retry without timeout to check this is done until finished
	worker2 := createWorker()
	worker2.Run(0)
	time.Sleep(250 * time.Millisecond)
	if state := worker2.GetState(); state != gotask.Finished {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(state))
	}
}

func TestReset(t *testing.T) {

	worker := createWorker()
	if state := worker.GetState(); state != gotask.Waiting {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Waiting), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 6 {
		t.Errorf("remaining weight not equal to 6: %v", weight)
	}

	worker.Run(0)
	worker.Wait()
	if state := worker.GetState(); state != gotask.Finished {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Finished), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 0 {
		t.Errorf("remaining weight not equal to 0: %v", weight)
	}

	worker.Reset()
	if state := worker.GetState(); state != gotask.Waiting {
		t.Errorf("worker state not equal to %v: %v", gotask.StateToString(gotask.Waiting), gotask.StateToString(state))
	}
	if weight := worker.GetRemainingWorkLoad(); weight != 6 {
		t.Errorf("remaining weight not equal to 6: %v", weight)
	}
}

func TestGetProgress(t *testing.T) {

	worker := createWorker()
	prog := worker.GetProgress()
	if prog != gotask.MinProgress {
		t.Errorf("remaining progress not %v, got: %v", gotask.MinProgress, prog)
	}

	worker.Run(0)
	time.Sleep(55 * time.Millisecond) // first task over
	prog = worker.GetProgress()
	if prog < 16.6 || prog > 17 {
		t.Errorf("remaining progress not %v, got: %v", 16, prog)
	}

	time.Sleep(50 * time.Millisecond) // second task over
	prog = worker.GetProgress()
	if prog != 50 {
		t.Errorf("remaining progress not %v, got: %v", 0.5, prog)
	}

	time.Sleep(50 * time.Millisecond) // last task over
	prog = worker.GetProgress()
	if prog != gotask.MaxProgress {
		t.Errorf("remaining progress not %v, got: %v", gotask.MaxProgress, prog)
	}

}
