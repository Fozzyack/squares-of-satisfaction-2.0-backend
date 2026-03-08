-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT users_email_key;
-- +goose StatementEnd
