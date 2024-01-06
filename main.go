package main

import (
	"math/rand"
	"task_flower/pkg/task"
	"task_flower/pkg/task_flow"
	"time"
)

func main() {
	tasksId := []string{"a1", "a2", "a3", "a4", "b1", "b2", "c1", "c2", "c3", "c4", "d1"}
	tasks := make(map[string]*task.Task)

	timeToSleep := []time.Duration{0 * time.Second, 10 * time.Second, 3 * time.Second, 5 * time.Second}

	for _, taskId := range tasksId {
		randomIndex := rand.Intn(len(timeToSleep))
		tasks[taskId] = task.New(taskId).Print(taskId).Sleep(timeToSleep[randomIndex]).Build()
	}

	constrains := map[string][]string{"a1": {"b1"}, "a2": {"b1"}, "a3": {"b2"}, "a4": {}, "b1": {"c1", "c2", "c3"}, "b2": {"c2", "c3", "c4"}, "c1": {"d1"}}

	taskFlower := task_flow.NewTaskFlow(tasks, constrains)
	_ = taskFlower.Run()
}
