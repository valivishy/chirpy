-- +goose Up
CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    email      TEXT      NOT NULL
);

CREATE UNIQUE INDEX idx_uk_users_email ON users (email);

-- +goose Down
DROP TABLE users;