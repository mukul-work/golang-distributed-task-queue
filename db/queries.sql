-- name: CreateJob :one
INSERT INTO jobs (task, max_attempts) VALUES ($1, $2) RETURNING *;

-- name: GetPendingJobs :many
SELECT * FROM jobs WHERE status = 'pending' ORDER BY created_at LIMIT $1;

-- name: MarkProcessing :one
UPDATE jobs SET status = 'processing', updated_at = now() WHERE id = $1 RETURNING *;

-- name: MarkDone :exec
UPDATE jobs SET status = 'done', updated_at = now() WHERE id = $1;

-- name: FailOrRetry :one
UPDATE jobs
SET attempts = attempts + 1,
    status = CASE WHEN attempts + 1 >= max_attempts THEN 'failed' ELSE 'pending' END,
    updated_at = now()
WHERE id = $1
RETURNING *;