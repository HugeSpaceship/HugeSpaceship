INSERT INTO migrations (migration, succeeded)
VALUES ('006_Location_composite', false);

BEGIN TRANSACTION;

CREATE TYPE earth_location AS (
    loc_X int,
    loc_Y int
);

UPDATE migrations
SET succeeded = true
WHERE migration = '006_Location_composite';
COMMIT;