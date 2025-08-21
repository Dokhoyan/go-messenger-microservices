-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    birth_date TIMESTAMP NOT NULL,
    avatar_url TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;

