-- +goose Up
CREATE TABLE refresh_tokens
(
    token      TEXT PRIMARY KEY NOT NULL,
    created_at TIMESTAMP        NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP        NOT NULL DEFAULT NOW(),
    user_id    UUID             NOT NULL,
    expires_at TIMESTAMP        NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMP,
    CONSTRAINT fk_refresh_tokens_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;