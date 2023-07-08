BEGIN TRANSACTION;

INSERT INTO migrations (migration, succeeded)
VALUES ('003_Slots', false);

ALTER TABLE users ADD COLUMN psn_uid text UNIQUE NOT NULL default '0';
ALTER TABLE users ADD COLUMN rpcn_uid text UNIQUE NOT NULL default '0';

CREATE TABLE slots (
    id uuid PRIMARY KEY,
    name text,
    description text,
    icon        char[40] NOT NULL REFERENCES resources(hash),
    root_level  char[40] NOT NULL REFERENCES resources(hash),
    locationX int4,
    locationY int4,
    initially_locked bool,
    sub_level bool,
    lbp1only bool,
    shareable int,
    background text,
    level_type                   text,
    min_players                  int,
    max_players                  int,
    move_required                bool,
    first_published timestamp
);

CREATE TABLE slot_resources (
  slot_id uuid REFERENCES slots(ID),
  resource_hash char[40] REFERENCES resources(hash)
);

UPDATE migrations
SET succeeded = true
WHERE migration = '003_Slots';

COMMIT;