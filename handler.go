package gotask

type State uint8      // State of Task and or worker
type Progress float64 // Task Progress
type Weight float64   // Weighting of task; a weight of 1 resembles ca. 1 second work time

const MinProgress Progress = 0   // minimum progress value in percent
const MaxProgress Progress = 100 // maximum progress value in percent

const (
	Waiting  State = iota // Task or Worker are ready to start
	Running  State = iota // Task or Worker currently running
	Canceled State = iota // Worker was stopped before finished due to timeout or due to user cancelled it
	Finished State = iota // Task or Worker finished. To rerun again call the reset method
)

// StateToString Converts task state to string equivalent
func StateToString(state State) string {
	switch state {
	case Waiting:
		return "WAITING"
	case Running:
		return "RUNNING"
	case Canceled:
		return "CANCELED"
	case Finished:
		return "FINISHED"
	}
	return ""
}
