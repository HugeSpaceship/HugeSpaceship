DROP TABLE plays;

CREATE TABLE scoreboard (
    id uuid primary key,
    slot_id int references slots,
    achieved_time timestamptz,

    score bigint,
    type int,
    players varchar[4] not null,
    main_player uuid references users,
    platform platform not null
);