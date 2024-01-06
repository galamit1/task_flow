package task

import "fmt"

type PrintTask struct {
	out string
}

func (t PrintTask) Run() error {
	fmt.Printf("%s\n", t.out)
	return nil
}

func NewPrintTask(id, out string) Task {
	return Task{
		ID:   id,
		Func: PrintTask{out: out},
	}
}
