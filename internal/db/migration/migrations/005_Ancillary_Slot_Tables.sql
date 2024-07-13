CREATE TABLE hearts (
    slot_id int references slots(id),
    owner uuid references users(id),
    heart_time timestamp
);

CREATE TABLE thumbs_up (
    slot_id int references slots(id),
    owner uuid references users(id),
    thumb_up_time timestamp
);
CREATE TABLE thumbs_down (
    slot_id int references slots(id),
    owner uuid references users(id),
    thumb_down_time timestamp
);

CREATE TABLE plays (
    slot_id int references slots(id),
    owner uuid references users(id),
    play_time timestamp,
    players int,
    score bigint
);