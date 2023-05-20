-- Migrations table stores what migrations we've done as to not trip over ourselves
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS hs_migrations (
    id serial,
    migration varchar(255),
    succeeded bool
);

INSERT INTO hs_migrations (migration) VALUES ('000_Initial_Migration');

COMMIT;

BEGIN TRANSACTION;

CREATE TABLE hs_users (
    id serial primary key,
    username varchar(20),
    avatar_hash bytea
);

CREATE TABLE hs_slots (
    id serial primary key,
    uploader integer references hs_users(id),

);

COMMIT;