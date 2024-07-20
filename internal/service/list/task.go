package list

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"todo-list/domain/task"
	"todo-list/log"
)

type Service struct {
	repo task.Repository
}

func New(opts ...func(s *Service) error) Service {
	s := Service{}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}

func WithTaskRepository(taskRepository task.Repository) func(s *Service) error {
	return func(s *Service) error {
		s.repo = taskRepository
		return nil
	}
}

func (s *Service) CreateTask(ctx context.Context, i task.Request) (id string, err error) {
	logger := log.LoggerFromContext(ctx)

	id = generateHex24()

	data := task.Entity{
		ID:       id,
		Title:    i.Title,
		ActiveAt: task.OnlyDate(i.ActiveAt),
		Status:   i.Status,
	}

	if err = s.repo.Create(ctx, data); err != nil {
		logger.Err(err).Stack().Msg("failed to create")
		return
	}

	return
}

func (s *Service) UpdateTask(ctx context.Context, id string, i task.Request) (err error) {
	logger := log.LoggerFromContext(ctx)

	data := task.Entity{
		Title:    i.Title,
		ActiveAt: task.OnlyDate(i.ActiveAt),
		Status:   i.Status,
	}

	if err = s.repo.Update(ctx, id, data); err != nil {
		logger.Err(err).Stack().Msg("failed to update")
		return
	}

	return
}

func (s *Service) DeleteTask(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx)

	if err = s.repo.Delete(ctx, id); err != nil {
		logger.Err(err).Stack().Msg("failed to delete")
		return
	}

	return
}

func (s *Service) DoneTask(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx)

	if err = s.repo.Done(ctx, id); err != nil {
		logger.Err(err).Stack().Msg("failed to done")
		return
	}

	return
}

func (s *Service) ListTasks(ctx context.Context, status string) (tasks []task.Entity, err error) {
	logger := log.LoggerFromContext(ctx)

	if status == "" {
		status = "active"
	}

	if status == "active" {
		tasks, err = s.repo.ListActive(ctx)
		if err != nil {
			logger.Err(err).Stack().Msg("failed to list active")
			return
		}
	} else if status == "done" {
		tasks, err = s.repo.ListDone(ctx)
		if err != nil {
			logger.Err(err).Stack().Msg("failed to list done")
			return
		}
	}

	for i := range tasks {
		_, err = tasks[i].TransformByWeekday()

		if err != nil {
			logger.Error().Err(err).Msg("failed to transform task")
			return
		}
	}

	return
}

func generateHex24() string {
	bytes := make([]byte, 12)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
