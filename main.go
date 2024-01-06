package main

import (
	"task_flower/pkg/task"
	"task_flower/pkg/task_flow"
)

func main() {
	tasksId := []string{"a1", "a2", "a3", "a4", "b1", "b2", "c1", "c2", "c3", "c4", "d1"}
	tasks := make(map[string]task.Task)
	
	for _, taskId := range tasksId {
		tasks[taskId] = task.NewPrintTask(taskId, taskId)
	}

	constrains := map[string][]string{"a1": {"b1"}, "a2": {"b1"}, "a3": {"b2"}, "a4": {}, "b1": {"c1", "c2", "c3"}, "b2": {"c2", "c3", "c4"}, "c1": {"d1"}}

	taskFlower := task_flow.NewTaskFlow(tasks, constrains)
	_ = taskFlower.Run()
}
