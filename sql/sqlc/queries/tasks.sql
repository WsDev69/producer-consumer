-- name: CreateTask :one
INSERT INTO tasks (type, value, creation_time)
VALUES ($1, $2, NOW())
RETURNING id, type, value, state, creation_time, last_update_time;

-- name: GetTaskByID :one
SELECT id, type, value, state, creation_time, last_update_time
FROM tasks
WHERE id = $1;

-- name: ListTasks :many
SELECT id, type, value, state, creation_time, last_update_time
FROM tasks
ORDER BY creation_time DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTaskState :exec
UPDATE tasks
SET state = $2
WHERE id = $1;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;

-- name: GetUnprocessedCount :one
SELECT COUNT(id) FROM tasks
WHERE state = 'received';
