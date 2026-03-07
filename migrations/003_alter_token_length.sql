-- +goose Up
-- +goose StatementBegin
ALTER TABLE sessions ALTER COLUMN token TYPE VARCHAR(511) ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions ALTER COLUMN token TYPE VARCHAR(63);
-- +goose StatementEnd
