-- +goose Up
ALTER TABLE transfers ADD COLUMN idempotency_key VARCHAR(255) UNIQUE;

-- +goose Down
ALTER TABLE transfers DROP COLUMN idempotency_key;
