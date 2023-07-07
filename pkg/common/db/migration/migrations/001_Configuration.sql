
BEGIN TRANSACTION;

INSERT INTO migrations (migration, succeeded)
VALUES ('001_Configuration', false);


CREATE TABLE config
(
    section text,
    config json
);

UPDATE migrations
SET succeeded = true
WHERE migration = '001_Configuration';

COMMIT;
