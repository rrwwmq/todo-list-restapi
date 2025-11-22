package todo

import "context"

type TaskRepository interface {
	AddTask(ctx context.Context, t Task) error
	GetTask(ctx context.Context, title string) (Task, error)
	ListTasks(ctx context.Context) ([]Task, error)
	CompleteTask(ctx context.Context, title string, completed bool) error
	UnCompleteTask(ctx context.Context, title string, completed bool) error
	ListUncompletedTasks(ctx context.Context) ([]Task, error)
	DeleteTask(ctx context.Context, title string) error
}
