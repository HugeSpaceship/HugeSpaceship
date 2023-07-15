-- Migrations table stores what migrations we've done as to not trip over ourselves
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS migrations
(
    id        serial primary key  not null,
    migration varchar(255) unique not null,
    succeeded bool                not null
);

COMMIT;

BEGIN TRANSACTION;

INSERT INTO migrations (migration, succeeded)
VALUES ('000_Initial_Migration', false);

CREATE TABLE users
(
    id          uuid primary key,
    username    varchar(20),
    avatar_hash bytea
);

CREATE TYPE platform AS ENUM ( -- what platform are they on,
    'PS3',
    'PSP',
    'PSVita',
    'PS4',
    'RPCS3',
    'Browser'
    );

CREATE TYPE game AS ENUM ( -- What game the session is for, or the website
    'LBP1',
    'LBPPSP',
    'LBP2',
    'LBPV',
    'LBP3',
    'Web'
    );

CREATE TABLE sessions
(
    id       serial primary key,
    userId   uuid references users,
    ip       inet     not null,
    token    uuid     not null,
    game     game     not null,
    platform platform not null
);

UPDATE migrations
SET succeeded = true
WHERE migration = '000_Initial_Migration';

COMMIT;