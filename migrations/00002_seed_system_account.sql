-- +goose Up
INSERT INTO accounts_types(name) VALUES ('system'), ('personal');

INSERT INTO accounts (account_type_id, is_system)
SELECT id, true FROM accounts_types WHERE name = 'system';
-- +goose Down
DELETE FROM accounts WHERE is_system;

DELETE FROM accounts_types WHERE name IN ('system', 'personal');