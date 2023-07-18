BEGIN TRANSACTION;

INSERT INTO migrations (migration, succeeded)
VALUES ('004_Slot_changes', false);

DROP TABLE slots CASCADE;
DROP TABLE slot_resources;

CREATE TABLE slots
(
    id               serial PRIMARY KEY,

    uploader         uuid REFERENCES users(id),

    name             text,
    description      text,
    icon             char(40) NOT NULL, -- Icons can be blank, therefore the FK doesn't help
    root_level       char(40) NOT NULL REFERENCES resources (hash),
    locationX        int4,
    locationY        int4,
    initially_locked bool,
    sub_level        bool,
    lbp1only         bool,
    shareable        int,
    background       text,
    level_type       text,
    min_players      int,
    max_players      int,
    move_required    bool,
    first_published  timestamp,
    last_updated     timestamp,
    domain           int

);

ALTER SEQUENCE slots_id_seq RESTART WITH 1;

CREATE TABLE slot_resources (
                                slot_id int REFERENCES slots(ID),
                                resource_hash char(40) REFERENCES resources(hash)
);

UPDATE migrations
SET succeeded = true
WHERE migration = '004_Slot_changes';

COMMIT;