package task

import (
	"fmt"
	"time"
)

type TaskBuilder struct {
	id          string
	out         string
	timeToSleep *time.Duration
}

func New(id string) *TaskBuilder {
	return &TaskBuilder{
		id: id,
	}
}

func (builder *TaskBuilder) Print(out string) *TaskBuilder {
	builder.out = out
	return builder
}

func (builder *TaskBuilder) Sleep(duration time.Duration) *TaskBuilder {
	builder.timeToSleep = &duration
	return builder
}

func (builder *TaskBuilder) Build() *Task {
	return &Task{
		ID: builder.id,
		Func: func() error {
			fmt.Printf("%v\n", builder.out)

			if builder.timeToSleep != nil {
				time.Sleep(*builder.timeToSleep)
				fmt.Printf("finish sleeping %v for task %s\n", builder.timeToSleep, builder.id)
			}

			return nil
		},
	}
}
