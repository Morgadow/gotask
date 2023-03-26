package gotask

type State uint8      // State of Task and or worker
type Progress float64 // Worker and Task Progress, Note: Task can be either done or not done, it has no float progress
type Weight float64   // Weighting of task; a weight of 1 resembles ca. 1 second work time

const MinProgress Progress = 0   // minimum progress value in percent
const MaxProgress Progress = 100 // maximum progress value in percent

const (
	Waiting        State = iota // Task or Worker are ready to start
	Running        State = iota // Task or Worker currently running
	Canceled       State = iota // Worker was stopped before finished due to timeout or due to user cancelled it
	Finished       State = iota // Task or Worker finished. To rerun again call the reset method
	TimeoutReached State = iota // Worker did not finish in time, equal to Canceled
)

var stateToString = map[State]string{Waiting: "WAITING", Running: "RUNNING", Canceled: "Canceled", Finished: "FINISHED", TimeoutReached: "TIMEOUT"}
var stringToState = map[string]State{"WAITING": Waiting, "RUNNING": Running, "CANCELED": Canceled, "FINISHED": Finished, "TIMEOUT": TimeoutReached}

// StateToString Converts task state to string equivalent
func StateToString(state State) string {
	return stateToString[state]
}

// StateToString Converts task state to string equivalent
func StringToState(name string) State {
	return stringToState[name]
}
