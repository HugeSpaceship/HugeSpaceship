ALTER TABLE users ADD COLUMN bio text NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN comments_enabled bool not null default true ;
ALTER TABLE users ADD COLUMN locationX integer not null default 0;
ALTER TABLE users ADD COLUMN locationY integer not null default 0;
alter table users
    alter column avatar_hash type text using avatar_hash::text;
