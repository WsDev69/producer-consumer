BEGIN;

-- Drop the trigger first to avoid dependency issues
DROP TRIGGER IF EXISTS update_task_last_update_time ON tasks;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_last_update_time();

-- Drop the tasks table
DROP TABLE IF EXISTS tasks;

-- Drop the TASK_STATE enum type
DROP TYPE IF EXISTS TASK_STATE;

COMMIT;