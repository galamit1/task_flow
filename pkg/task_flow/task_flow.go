package task_flow

import (
	"fmt"
	"sync"
	"task_flower/pkg/task"
)

type TaskFlow struct {
	taskIDToTask      map[string]*task.Task
	constrains        map[string][]string
	runningTasks      map[string]struct{}
	finishedTasksChan chan string
	// a wait-group to make sure all running tasks are finished
	waitGroup *sync.WaitGroup
}

func NewTaskFlow(taskIDToTask map[string]*task.Task, constrains map[string][]string) TaskFlow {
	return TaskFlow{
		taskIDToTask:      taskIDToTask,
		constrains:        constrains,
		runningTasks:      make(map[string]struct{}),
		finishedTasksChan: make(chan string),
		waitGroup:         &sync.WaitGroup{},
	}
}

func (flow TaskFlow) getHeadTasks() []string {
	candidatesForHead := make(map[string]struct{})
	// add all the blocking tasks in the constrains to be the candidates
	for id, _ := range flow.taskIDToTask {
		candidatesForHead[id] = struct{}{}
	}

	// remove all the tasks that are blocked from candidatesForHead
	for blocking, tasks := range flow.constrains {
		for _, taskToRemove := range tasks {
			if _, ok := flow.taskIDToTask[blocking]; ok { // removes only if the block task exists
				delete(candidatesForHead, taskToRemove)
			}
		}
	}

	// remove all the tasks that are currently running
	for taskId, _ := range flow.runningTasks {
		delete(candidatesForHead, taskId)
	}

	headTasks := make([]string, 0, len(candidatesForHead))
	for headTask, _ := range candidatesForHead {
		headTasks = append(headTasks, headTask)
	}

	return headTasks
}

func (flow TaskFlow) runAvailableTasks() {
	headTasks := flow.getHeadTasks()
	if len(headTasks) == 0 {
		return
	}

	fmt.Printf("running tasks %v\n", headTasks)
	err := flow.runTasks(headTasks)
	if err != nil {
		fmt.Printf("error running tasks %v, err: %v", headTasks, err)
	}
}

func (flow TaskFlow) runTasks(tasks []string) error {
	for _, taskId := range tasks {
		flow.waitGroup.Add(1)
		flow.runningTasks[taskId] = struct{}{}

		go func(taskId string) {
			defer func() {
				flow.waitGroup.Done()
				flow.finishedTasksChan <- taskId
			}()
			function, ok := flow.taskIDToTask[taskId]
			if ok {
				function.Func()
			}
		}(taskId)
	}

	return nil

}

func (flow TaskFlow) Run() error {
	flow.runAvailableTasks()

	for taskToRemove := range flow.finishedTasksChan {
		// remove the finished task so we can run the next tasks
		delete(flow.taskIDToTask, taskToRemove)
		delete(flow.runningTasks, taskToRemove)

		flow.runAvailableTasks()

		// check if we finished run all the tasks
		if len(flow.taskIDToTask) == 0 {
			break
		}
	}

	// wait for all the tasks to finish
	flow.waitGroup.Wait()
	close(flow.finishedTasksChan)

	return nil
}
