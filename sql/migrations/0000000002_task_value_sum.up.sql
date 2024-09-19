BEGIN;

-- Create the task_value_sum table
CREATE TABLE task_sums
(
    task_type        INTEGER PRIMARY KEY,       -- The type (between 0 and 9)
    total_value INTEGER NOT NULL DEFAULT 0 -- The sum of values for this type
);


CREATE OR REPLACE FUNCTION update_sum_on_insert() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.state = 'done' THEN
        INSERT INTO task_sums (task_type, total_value)
        VALUES (NEW.type, NEW.value)
        ON CONFLICT (task_type) DO UPDATE
        SET total_value = task_sums.total_value + EXCLUDED.total_value;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_sum_insert
    AFTER INSERT ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_sum_on_insert();

CREATE OR REPLACE FUNCTION update_sum_on_update() RETURNS TRIGGER AS $$
BEGIN
    -- If the state is changing to 'done'
    IF OLD.state != 'done' AND NEW.state = 'done' THEN
        -- Add the new value to the sum
        INSERT INTO task_sums (task_type, total_value)
        VALUES (NEW.type, NEW.value)
        ON CONFLICT (task_type) DO UPDATE
                                         SET total_value = task_sums.total_value + EXCLUDED.total_value;

-- If the state was 'done' and is changing to something else
ELSIF OLD.state = 'done' AND NEW.state != 'done' THEN
        -- Subtract the old value from the sum
UPDATE task_sums
SET total_value = total_value - OLD.value
WHERE task_type = OLD.type;

-- If the state remains 'done', but the value is changing
ELSIF OLD.state = 'done' AND NEW.state = 'done' THEN
        -- Adjust the sum by subtracting the old value and adding the new one
UPDATE task_sums
SET total_value = total_value - OLD.value + NEW.value
WHERE task_type = OLD.type;
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_sum_update
    AFTER UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_sum_on_update();


CREATE OR REPLACE FUNCTION update_sum_on_delete() RETURNS TRIGGER AS $$
BEGIN
    IF OLD.state = 'done' THEN
UPDATE task_sums
SET total_value = total_value - OLD.value
WHERE task_type = OLD.type;
END IF;
RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_sum_delete
    AFTER DELETE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_sum_on_delete();


COMMIT;
