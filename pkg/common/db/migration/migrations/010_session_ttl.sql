INSERT INTO migrations (migration, succeeded)
VALUES ('010_session_ttl', false);

BEGIN TRANSACTION;

ALTER TABLE sessions ADD COLUMN expiry timestamp without time zone;

UPDATE migrations
SET succeeded = true
WHERE migration = '010_session_ttl';
COMMIT;