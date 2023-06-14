
BEGIN TRANSACTION;

INSERT INTO migrations (migration, succeeded)
VALUES ('001_Configuration', false);

CREATE EXTENSION "hstore";

CREATE TABLE config
(
    section text,
    values hstore
);

UPDATE migrations
SET succeeded = true
WHERE migration = '001_Configuration';

COMMIT;
