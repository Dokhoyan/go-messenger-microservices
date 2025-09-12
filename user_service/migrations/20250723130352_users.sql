-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    birth_date TIMESTAMP NOT NULL,
    avatar_url TEXT NOT NULL,
    role INTEGER,
    password VARCHAR,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
DROP TABLE IF EXISTS users;

