package task

type TaskFunc interface {
	Run() error
}

type Task struct {
	ID   string
	Func TaskFunc
}
