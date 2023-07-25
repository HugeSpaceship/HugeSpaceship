INSERT INTO migrations (migration, succeeded)
VALUES ('008_User_changes', false);

BEGIN TRANSACTION;

ALTER TABLE users ADD COLUMN bio text NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN comments_enabled bool not null default true ;
ALTER TABLE users ADD COLUMN locationX integer not null default 0;
ALTER TABLE users ADD COLUMN locationY integer not null default 0;
alter table users
    alter column avatar_hash type text using avatar_hash::text;


UPDATE migrations
SET succeeded = true
WHERE migration = '008_User_changes'; -- I removed this because it would be annoying to implement
COMMIT;