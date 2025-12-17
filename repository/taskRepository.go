package repository

import (
	"context"
	"errors"
	"time"

	"firstRestAPI/database"
	"firstRestAPI/todo"
)

var ErrTaskNotFound = errors.New("task not found")

type GormRepository struct{}

func NewGormRepository() *GormRepository {
	return &GormRepository{}
}

func (r *GormRepository) AddTask(ctx context.Context, t todo.Task) error {
	return database.DB.WithContext(ctx).Create(&t).Error
}

func (r *GormRepository) GetTask(ctx context.Context, title string) (todo.Task, error) {
	var t todo.Task
	err := database.DB.WithContext(ctx).Where("title = ?", title).First(&t).Error
	if err != nil {
		return todo.Task{}, ErrTaskNotFound
	}
	return t, nil
}

func (r *GormRepository) ListTasks(ctx context.Context) ([]todo.Task, error) {
	var tasks []todo.Task
	err := database.DB.WithContext(ctx).Find(&tasks).Error
	return tasks, err
}

func (r *GormRepository) CompleteTask(ctx context.Context, title string, completed bool) error {
	var completedAt *time.Time
	if completed {
		now := time.Now()
		completedAt = &now
	}
	return database.DB.WithContext(ctx).Model(&todo.Task{}).
		Where("title = ?", title).
		Updates(map[string]interface{}{
			"is_completed": completed,
			"completed_at": completedAt,
		}).Error
}

func (r *GormRepository) UnCompleteTask(ctx context.Context, title string, completed bool) error {
	var completedAt *time.Time
	if !completed {
		completedAt = nil
	}
	return database.DB.WithContext(ctx).Model(&todo.Task{}).
		Where("title = ?", title).
		Updates(map[string]interface{}{
			"is_completed": completed,
			"completed_at": completedAt,
		}).Error
}

func (r *GormRepository) ListUncompletedTasks(ctx context.Context) ([]todo.Task, error) {
	var tasks []todo.Task
	err := database.DB.WithContext(ctx).
		Where("is_completed = ?", false).
		Order("created_at desc").
		Find(&tasks).Error
	return tasks, err
}

func (r *GormRepository) DeleteTask(ctx context.Context, title string) error {
	return database.DB.WithContext(ctx).
		Where("title = ?", title).
		Delete(&todo.Task{}).Error
}
