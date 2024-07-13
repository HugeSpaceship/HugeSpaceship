CREATE TABLE IF NOT EXISTS users
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

CREATE TABLE IF NOT EXISTS sessions
(
    id       serial primary key,
    userId   uuid references users,
    ip       inet     not null,
    token    uuid     not null,
    game     game     not null,
    platform platform not null
);