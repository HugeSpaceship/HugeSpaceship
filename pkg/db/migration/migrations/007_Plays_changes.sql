ALTER TABLE plays ADD COLUMN game_type game NOT NULL DEFAULT 'LBP2';

ALTER TABLE slots ALTER first_published SET DEFAULT 'epoch';
ALTER TABLE slots ALTER first_published SET NOT NULL;

ALTER TABLE slots ALTER last_updated SET DEFAULT 'epoch';
ALTER TABLE slots ALTER last_updated SET NOT NULL;

ALTER TABLE users ADD COLUMN entitled_slots integer NOT NULL DEFAULT 50;