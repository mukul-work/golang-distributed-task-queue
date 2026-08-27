package service

import (
	"context"

	"github.com/mukul-work/golang-distributed-task-queue/internal/dbgen"
	"github.com/mukul-work/golang-distributed-task-queue/internal/store"
)

type Queue interface {
	Enqueue(ctx context.Context, job dbgen.Job) error
	Dequeue(ctx context.Context) (*dbgen.Job, error)
	GetJob(ctx context.Context, id string) (dbgen.Job, error)
}

type Service struct {
	s *store.Store
}

func NewService(s *store.Store) *Service {
	if s == nil {
		return nil
	}
	return &Service{
		s: s,
	}
}

func (s *Service) Enqueue(ctx context.Context, job dbgen.Job) error {
	arg := dbgen.CreateJobParams{
		ID:          job.ID,
		Type:        job.Type,
		Payload:     job.Payload,
		Status:      job.Status,
		MaxAttempts: job.MaxAttempts,
	}
	_, err := s.s.CreateJob(ctx, arg)
	return err
}

func (s *Service) Dequeue(ctx context.Context) (*dbgen.Job, error) {
	job, err := s.s.DequeueJob(ctx)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
