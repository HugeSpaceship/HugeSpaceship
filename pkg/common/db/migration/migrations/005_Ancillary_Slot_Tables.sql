BEGIN TRANSACTION;

INSERT INTO migrations (migration, succeeded)
VALUES ('005_Ancillary_Slot_Tables', false);

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



UPDATE migrations
SET succeeded = true
WHERE migration = '004_Slot_changes';
COMMIT;