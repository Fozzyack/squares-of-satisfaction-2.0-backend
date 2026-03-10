-- +goose Up
-- +goose StatementBegin
ALTER TABLE habit_entries ADD COLUMN date DATE;

UPDATE habit_entries
SET date = (created_at AT TIME ZONE 'UTC')::date
WHERE date IS NULL;

ALTER TABLE habit_entries
ALTER COLUMN date SET NOT NULL,
ALTER COLUMN date SET DEFAULT CURRENT_DATE;

CREATE INDEX idx_habit_entries_user_id_date ON habit_entries (user_id, date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_habit_entries_user_id_date;
ALTER TABLE habit_entries DROP COLUMN date;
-- +goose StatementEnd
