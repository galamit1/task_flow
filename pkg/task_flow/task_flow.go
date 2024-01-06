package task_flow

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"task_flower/pkg/task"
)

type TaskFlow struct {
	taskIDToTask map[string]task.Task
	constrains   map[string][]string
}

func NewTaskFlow(taskIDToTask map[string]task.Task, constrains map[string][]string) TaskFlow {
	return TaskFlow{
		taskIDToTask: taskIDToTask,
		constrains:   constrains,
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

	headTasks := make([]string, 0, len(candidatesForHead))
	for headTask, _ := range candidatesForHead {
		headTasks = append(headTasks, headTask)
	}

	return headTasks
}

func (flow TaskFlow) removeTasks(tasks []string) {
	for _, taskToRemove := range tasks {
		delete(flow.taskIDToTask, taskToRemove)
	}
}

func (flow TaskFlow) runTasks(tasks []string) error {
	eg := errgroup.Group{}
	for _, taskId := range tasks {
		fmt.Printf("run func %v\n", taskId)
		eg.Go(func() error {
			function, ok := flow.taskIDToTask[taskId]
			if ok {
				return function.Func.Run()
			}

			return nil
		})
	}

	return eg.Wait()
}

func (flow TaskFlow) Run() error {
	for len(flow.taskIDToTask) > 0 {
		headTasks := flow.getHeadTasks()
		fmt.Printf("running tasks %v\n", headTasks)
		err := flow.runTasks(headTasks)
		if err != nil {
			fmt.Printf("error running tasks %v, err: %v", headTasks, err)
		}
		flow.removeTasks(headTasks)
	}

	return nil
}
