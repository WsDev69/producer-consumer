BEGIN;

DROP TRIGGER IF EXISTS trigger_sum_insert ON tasks;
DROP TRIGGER IF EXISTS trigger_sum_update ON tasks;
DROP TRIGGER IF EXISTS trigger_sum_delete ON tasks;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_sum_on_insert();
DROP FUNCTION IF EXISTS update_sum_on_update();
DROP FUNCTION IF EXISTS update_sum_on_delete();

-- Drop the tasks table
DROP TABLE IF EXISTS task_sums;


COMMIT;