INSERT INTO migrations (migration, succeeded)
VALUES ('007_Location_removal', false);

BEGIN TRANSACTION;

UPDATE migrations
SET succeeded = true
WHERE migration = '007_Location_removal'; -- I removed this because it would be annoying to implement
COMMIT;