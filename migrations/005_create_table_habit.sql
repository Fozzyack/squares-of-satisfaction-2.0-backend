-- +goose Up
-- +goose StatementBegin
CREATE TABLE habits (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    name VARCHAR(50) NOT NULL,
    goal INT NOT NULL,
    increment INT DEFAULT 1 NOT NULL, 
    color VARCHAR(31),

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE habits;
-- +goose StatementEnd
