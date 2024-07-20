package repository

import (
	"context"
	"errors"
	"todo-list/domain/task"

	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

type TaskRepository struct {
	db *sqlx.DB
}

func MustNew(logger *zerolog.Logger, db *sqlx.DB) *TaskRepository {
	if logger == nil {
		panic("logger is required")
	}

	if db == nil {
		panic("db is required")
	}

	return &TaskRepository{
		db: db,
	}
}

func (r TaskRepository) Create(ctx context.Context, t task.Entity) (err error) {
	q := `
		INSERT INTO tasks (id, title, active_at, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err = r.db.ExecContext(ctx, q, t.ID, t.Title, t.ActiveAt, t.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = task.ErrExists
		}
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return task.ErrExists
		}
		return
	}

	return
}

func (r *TaskRepository) Update(ctx context.Context, id string, t task.Entity) (err error) {
	q := `
	UPDATE tasks
	SET title = $1, active_at = $2
	WHERE id = $3 RETURNING id
	`

	if err = r.db.QueryRowContext(ctx, q, t.Title, t.ActiveAt, id).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = task.ErrNotFound
			return
		}
	}

	return
}

func (r *TaskRepository) Delete(ctx context.Context, id string) (err error) {
	q := `
	DELETE FROM tasks WHERE id = $1 RETURNING id
	`

	if err = r.db.QueryRowContext(ctx, q, id).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = task.ErrNotFound
			return
		}
	}

	return
}

func (r *TaskRepository) Done(ctx context.Context, id string) (err error) {
	q := `
		UPDATE tasks
		SET active_at = CURRENT_DATE, status = 'done'
		WHERE id = $1 RETURNING id
	`

	if err = r.db.QueryRowContext(ctx, q, id).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = task.ErrNotFound
			return
		}
	}

	return
}

func (r *TaskRepository) ListActive(ctx context.Context) (tasks []task.Entity, err error) {
	tasks = []task.Entity{}

	q := `
		SELECT id, title, active_at, status
		FROM tasks
		WHERE status = 'active' AND active_at <= CURRENT_DATE
		ORDER BY active_at;
	`

	err = r.db.SelectContext(ctx, &tasks, q)

	return
}

func (r *TaskRepository) ListDone(ctx context.Context) (tasks []task.Entity, err error) {
	tasks = []task.Entity{}

	q := `
		SELECT id, title, active_at
		FROM tasks
		WHERE status = 'done'
		ORDER BY active_at;
	`

	err = r.db.SelectContext(ctx, &tasks, q)
	if err != nil {
		return
	}

	return
}
