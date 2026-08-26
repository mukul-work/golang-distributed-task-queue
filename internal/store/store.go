package job

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

func (s *Store) FailOrRetry(ctx context.Context, id int32) (dbgen.Job, error) {
	job, err := s.q.FailOrRetry(ctx, id)
	if err != nil {
		return dbgen.Job{}, fmt.Errorf("failed to execute query: %w", err)
	}
	return job, nil
}

func (s *Store) GetPendingJobs(ctx context.Context, limit int32) ([]dbgen.Job, error) {
	jobs, err := s.q.GetPendingJobs(ctx, limit)
	if err != nil {
		return []dbgen.Job{}, fmt.Errorf("failed to fetch pending jobs: %w", err)

	}
	return jobs, nil
}

func (s *Store) MarkDone(ctx context.Context, id int32) error {
	err := s.q.MarkDone(ctx, id)
	if err != nil {
		return fmt.Errorf("job not done: %w", err)
	}
	return nil
}

func (s *Store) MarkProcessing(ctx context.Context, id int32) (dbgen.Job, error) {
	job, err := s.q.MarkProcessing(ctx, id)
	if err != nil {
		return dbgen.Job{}, fmt.Errorf("job marked as processing: %w", err)
	}
	return job, nil
}
