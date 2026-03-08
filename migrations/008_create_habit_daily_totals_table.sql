-- +goose Up
-- +goose StatementBegin
CREATE TABLE habit_daily_totals (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    amount INT NOT NULL,
    habit_id UUID REFERENCES habits(id),
    user_id UUID REFERENCES users(id),
    date DATE DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_habit_daily_totals UNIQUE (habit_id, user_id, date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE habit_daily_totals;
-- +goose StatementEnd
