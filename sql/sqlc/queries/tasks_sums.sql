-- name: GetTotalValueByTaskType :one
SELECT task_type, total_value FROM task_sums
WHERE task_type = $1;