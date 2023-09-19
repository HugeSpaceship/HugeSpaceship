INSERT INTO migrations (migration, succeeded)
VALUES ('012_Photos', false);

BEGIN TRANSACTION;

CREATE TYPE slotType AS ENUM ( -- what platform are they on,
    'pod',
    'user',
    'developer',
    'moon',
    'remotemoon',
    'local'
    );

CREATE TABLE photos (
    id serial primary key,
    domain int not null,
    author uuid references users(id),
    small text not null,
    medium text not null,
    large text not null,
    plan text not null,
    slotType slotType not null,
    slotField text -- ID in the case of developer, root level in the case of pod, etc
);

CREATE TABLE photo_subjects (
    photo_id int references photos(id),
    user_id uuid,
    name text,
    x1 int,
    y1 int,
    x2 int,
    y2 int
);

UPDATE migrations
SET succeeded = true
WHERE migration = '012_Photos';
COMMIT;