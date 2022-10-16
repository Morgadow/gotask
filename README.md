# gotask

*Gotask* is a super lightweight task handler in Golang.

It was created to handle a big amount of small tasks inside a UI or console application.
It`s main purpose is to handle execution timeout, progress status and and progress description management.

---

## Data

The main structs are the **Worker** struct which handles the task execution as explained later.
A **Worker** consists of a lot of small **Tasks** which are created and then added to the **Worker**. Every **Task** implements the interface **Runnable**. It may be theoretically possible to use a custom defined task struct which also implements **Runnable**. For a proper execution it is recommended to stick with the proposed **Task** data type.

The **Worker** and every **Task** status are mainly defined over their state of type *State*, which can have following values:

```golang
const (
 Waiting  State = iota // Task or Worker are ready to start
 Running  State = iota // Task or Worker currently running
 Canceled State = iota // Worker was stopped before finished due to timeout or due to user cancelled it
 Finished State = iota // Task or Worker finished. To rerun again call the reset method
)
```

Information about how much work was already done is provided over the *Progress* value in percent.

## Usage

Create a **Worker** which and can be later filled with **Tasks**.
The **Worker** is created using the dedicated factory function for proper struct initialization.

```golang
worker := gotask.NewWorker("Workername")
```

Then the **Tasks** are created which should be executed from within the **Worker**. Every **Task** is created using the dedicated factory function for proper struct initialization.

While creating the **Task** *name* and *description* are defined for later status information purposes which **Task** is currently handled.
The task *weight* marks the amount of work which is done inside the **Task**. A *weight* of 1 stands for a duration of about one second of work time.
This information is later used to estimate the over percentage of **Worker** progress.

The *target* parameter is the function to be executed inside the **Task**. This function does not take any input parameter and does not output any. If any data is needed here, it must be gathered using a get() function or similar.

```golang
func NewTask(name string, weight Weight, desc string, target func()) *Task {}
```

Fill **Worker** with multiple tasks using the *AddTask()* method.
Multiple **Tasks** can also be added at once using the *AddTasks()* method.
Once added, all subtasks can be deleted from the internal queue over the *ClearTasks()* method.

```golang
// Sleep One function to execute inside a task.
func Sleep() {
 time.Sleep(time.Millisecond * 250)
}

_ = worker.AddTask(gotask.NewTask("task 0", gotask.Weight(1),"sleeping for 250ms", Sleep))
_ = worker.AddTask(gotask.NewTask("task 1", gotask.Weight(1),"sleeping for 250ms", Sleep))
_ = worker.AddTask(gotask.NewTask("task 2", gotask.Weight(1),"sleeping for 250ms", Sleep))
_ = worker.AddTask(gotask.NewTask("task 3", gotask.Weight(1),"sleeping for 250ms", Sleep))
```

Once the **Worker** is filled with tasks it is ready to be executed. After starting the **Worker**, the main routing must wait until the **Worker** was finished before it should close itself.
> Note: Once the **Worker** is running, no more tasks can be added!

When starting the **Worker**, a *timeout* in seconds can be set as input parameter. To set no *timeout* at all, set a value equal or smaller to zero.

During worker runtime, it can be also stopped by the user over the *Stop()* method.
> Note: If a timeout is reached or the **Workers** *Stop()* method is called, the **Worker** does not stop immediatelly. It just does not start the next task in the queue and leaves with state **Canceled**! The Stop method therefore does not break endless loops inside the task target!

```golang
_ = worker.Run(5)  // start worker with a timeout of 5 seconds
err := worker.Wait() // main routine waits until worker is finished, a timeout is reached or stopped by user
fmt.Println("Worker finished with error: ", err)
```

Once the **Worker** is finished, it can be reset calling the *Reset()* method again. This method also resets the **Workers** state and progress and all the added tasks.

```golang
worker.Reset()
_ = worker.Run(0) // the worker is started without any timeout limitations
err = worker.Wait()
fmt.Println("Worker finished with error: ", err)
```

## Logging worker status during run

During the **Workers** runtime several informations can be requested. This example shows how, for example, a log mechanism can keep track of the **Workers** status.

```golang
// Function to print worker state and progress
func PrintWorkerStatus(worker *gotask.Worker) {
 for worker.GetState() == gotask.Running {
  prog := worker.GetProgress()
  wl := worker.GetRemainingWorkLoad()
  dur, _ := worker.GetDuration()
  remain, _ := worker.GetRemainingTime()
  currtaskname, _ := worker.GetCurrentTaskName()
  fmt.Printf("Progress: %v%%, Remaining Load: %v, Duration: %v s, Remain Time: %v s, CurrTask: %s\n", prog, wl, dur, remain, currtaskname)

  time.Sleep(480 * time.Millisecond)
 }

// Start worker and log method and then wait for the worker to finish
_ = worker.Run(5)  // if not timeout would be set, the remaining time is -1
go PrintWorkerStatus(worker)
err := worker.Wait()
fmt.Println("Worker finished with error: ", err)
}
```

This code will result in this log:

```golang
Progress: 0%, Remaining Load: 10, Duration: 0 s, Remain Time: 5 s, CurrTask: 
Progress: 10%, Remaining Load: 9, Duration: 0.48 s, Remain Time: 4.519 s, CurrTask: task 1
Progress: 30%, Remaining Load: 7, Duration: 0.96 s, Remain Time: 4.039 s, CurrTask: task 3
Progress: 50%, Remaining Load: 5, Duration: 1.441 s, Remain Time: 3.558 s, CurrTask: task 5
Progress: 70%, Remaining Load: 3, Duration: 1.921 s, Remain Time: 3.078 s, CurrTask: task 7
Progress: 90%, Remaining Load: 1, Duration: 2.401 s, Remain Time: 2.598 s, CurrTask: task 9
Worker finished with error:  <nil>
```

## Changelog

- **v1.0**: First working and tested release.
