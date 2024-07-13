CREATE TYPE nat_type AS ENUM (
    'open',
    'moderate',
    'strict'
    );

CREATE TYPE slot_type AS ENUM (
    'developer',
    'user',
    'moon',
    'pod',
    'local'
    );

CREATE TYPE room_slot AS (
    slot_id int,
    slot_type slot_type
);

CREATE TYPE room_user AS (
    player uuid,
    nat nat_type
);

CREATE TABLE rooms (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    players room_user[],
    game game NOT NULL,
    game_version text NOT NULL,
    platform platform NOT NULL,
    room_slot room_slot NOT NULL
);