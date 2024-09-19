BEGIN;

-- Create the enum type for task state
CREATE TYPE TASK_STATE AS ENUM (
    'received',    -- Initial state when the task is created
    'processing',  -- State when the task is being processed
    'done'         -- State when the task is completed
);

-- Create the tasks table
CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,           -- Auto-incrementing primary key
                       type INTEGER NOT NULL CHECK (type BETWEEN 0 AND 9),   -- integer between 0 and 9
                       value INTEGER NOT NULL CHECK (value BETWEEN 0 AND 99), -- integer between 0 and 99
                       state TASK_STATE NOT NULL DEFAULT 'received',   -- default is 'received'
                       creation_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),  -- Timestamp of task creation
                       last_update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()  -- Timestamp of the last update
);

CREATE INDEX idx_tasks_state ON tasks (state);

-- Add a trigger to update last_update_time on row update
CREATE OR REPLACE FUNCTION update_last_update_time()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_update_time = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to update last_update_time on updates
CREATE TRIGGER update_task_last_update_time
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_last_update_time();

COMMIT;
