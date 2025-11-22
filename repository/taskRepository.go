package repository

import (
	"context"
	"firstRestAPI/database"
	"firstRestAPI/todo"
)

type PostgresRepository struct{}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) AddTask(ctx context.Context, t todo.Task) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO Tasks (title, description)
		VALUES($1, $2)`,
		t.Title, t.Description,
	)

	return err
}

func (r *PostgresRepository) GetTask(ctx context.Context, title string) (todo.Task, error) {
	var task todo.Task

	err := database.DB.QueryRow(ctx,
		`SELECT
		id,
		title,
		description,
		is_completed,
		created_at,
		completed_at
	FROM Tasks
	WHERE title=$1`, title).Scan(&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.CompletedAt)

	if err != nil {
		return todo.Task{}, todo.ErrTaskNotFound
	}

	return task, nil
}

func (r *PostgresRepository) ListTasks(ctx context.Context) ([]todo.Task, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT
		id,
		title,
		description,
		is_completed,
		created_at,
		completed_at
	FROM Tasks
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []todo.Task

	for rows.Next() {
		var t todo.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.IsCompleted, &t.CreatedAt, &t.CompletedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *PostgresRepository) CompleteTask(ctx context.Context, title string, completed bool) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE Tasks
	SET is_completed = $1, completed_at = NOW()
	WHERE title = $2`, completed, title)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) UnCompleteTask(ctx context.Context, title string, completed bool) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE Tasks
	SET is_completed = $1, completed_at = NULL
	WHERE title = $2`, completed, title)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) ListUncompletedTasks(ctx context.Context) ([]todo.Task, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT
		id,
		title,
		description,
		is_completed,
		created_at,
		completed_at
	FROM Tasks
	WHERE is_completed = FALSE
	ORDER BY created_at DESC`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []todo.Task
	for rows.Next() {
		var t todo.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.IsCompleted, &t.CreatedAt, &t.CompletedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *PostgresRepository) DeleteTask(ctx context.Context, title string) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM Tasks WHERE title = $1`, title)

	if err != nil {
		return err
	}

	return nil
}
