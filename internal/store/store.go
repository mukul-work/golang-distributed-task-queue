package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mukul-work/golang-distributed-task-queue/internal/dbgen"
)

type Store struct {
	q     dbgen.Queries
	dbcon *pgxpool.Pool
}

func NewStore(dbcon *pgxpool.Pool) *Store {
	if dbcon == nil {
		return nil
	}

	return &Store{
		q:     *dbgen.New(dbcon),
		dbcon: dbcon,
	}
}

func (s *Store) CreateJob(ctx context.Context, arg dbgen.CreateJobParams) (dbgen.Job, error) {
	job, err := s.q.CreateJob(ctx, arg)
	if err != nil {
		return dbgen.Job{}, fmt.Errorf("unable to create Job: %w", err)

	}
	return job, nil
}

func (s *Store) GetJob(ctx context.Context, id string) (dbgen.Job, error) {
	job, err := s.q.GetJob(ctx, id)
	if err != nil {
		return dbgen.Job{}, fmt.Errorf("unable to fetch the job with the given id: %s", id)

	}
	return job, nil
}

func (s *Store) DequeueJob(ctx context.Context) (dbgen.Job, error) {
	job, err := s.q.DequeueJob(ctx)
	if err != nil {
		return dbgen.Job{}, fmt.Errorf("Could not dequeue job: %w", err)
	}
	return job, nil

}

func (s *Store) FailJob(ctx context.Context, arg dbgen.FailJobParams) error {
	return s.q.FailJob(ctx, arg)
}

func (s *Store) CompletedJob(ctx context.Context, arg dbgen.FailJobParams) error {
	return s.q.FailJob(ctx, arg)
}

func (s *Store) GetAPIKeyByHash(ctx context.Context, keyHash string) (dbgen.ApiKey, error) {
	return s.q.GetAPIKeyByHash(ctx, keyHash)
}

func (s *Store) UpdateLastUsed(ctx context.Context, id string) error {
	return s.q.UpdateLastUsed(ctx, id)
}
func (s *Store) CreateAPIKey(ctx context.Context, arg dbgen.CreateAPIKeyParams) (dbgen.ApiKey, error) {
	return s.q.CreateAPIKey(ctx, arg)
}
func (s *Store) ListAPIKeys(ctx context.Context) ([]dbgen.ApiKey, error) {
	return s.q.ListAPIKeys(ctx)
}
