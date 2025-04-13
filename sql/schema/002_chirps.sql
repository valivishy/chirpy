-- +goose Up
CREATE TABLE chirps
(
    id         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP        NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP        NOT NULL DEFAULT NOW(),
    body       TEXT             NOT NULL,
    user_id    UUID             NOT NULL,
    CONSTRAINT fk_chirps_users
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;