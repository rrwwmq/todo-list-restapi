package todo

import (
	"context"
)

type List struct {
	repo TaskRepository
}

func NewList(repo TaskRepository) *List {
	return &List{repo: repo}
}

func (l *List) AddTask(task Task) error {
	return l.repo.AddTask(context.Background(), task)
}

func (l *List) GetTask(title string) (Task, error) {
	return l.repo.GetTask(context.Background(), title)
}

func (l *List) ListTasks() ([]Task, error) {
	tasks, err := l.repo.ListTasks(context.Background())
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (l *List) CompleteTask(title string) (Task, error) {
	err := l.repo.CompleteTask(context.Background(), title, true)
	if err != nil {
		return Task{}, err
	}

	return l.repo.GetTask(context.Background(), title)
}

func (l *List) UncompleteTask(title string) (Task, error) {
	err := l.repo.UnCompleteTask(context.Background(), title, false)
	if err != nil {
		return Task{}, err
	}

	return l.repo.GetTask(context.Background(), title)
}

func (l *List) ListUncompletedTasks() ([]Task, error) {
	return l.repo.ListUncompletedTasks(context.Background())
}

func (l *List) DeleteTask(title string) error {
	return l.repo.DeleteTask(context.Background(), title)
}
