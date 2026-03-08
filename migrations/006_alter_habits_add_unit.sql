-- +goose Up
-- +goose StatementBegin
ALTER TABLE habits
ADD COLUMN unit VARCHAR(31);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE habits
DROP COLUMN unit;
-- +goose StatementEnd
