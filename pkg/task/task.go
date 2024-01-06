package task

type Task struct {
	ID   string
	Func func() error
}
