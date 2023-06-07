
BEGIN TRANSACTION;

INSERT INTO migrations (migration, succeeded)
VALUES ('001_Configuration', false);

CREATE EXTENSION "hstore";

CREATE TABLE config
(
    map hstore -- Apparently this is a key-value store
);

UPDATE migrations
SET succeeded = true
WHERE migration = '001_Configuration';

COMMIT;
