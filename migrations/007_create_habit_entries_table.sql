-- +goose Up
-- +goose StatementBegin
CREATE TABLE habit_entries (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    increment_amount INT NOT NULL,
    habit_id UUID REFERENCES habits(id),
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE habit_entries;
-- +goose StatementEnd
