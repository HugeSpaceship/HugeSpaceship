INSERT INTO migrations (migration, succeeded)
VALUES ('011_slot_resources_FK_fix', false);

BEGIN TRANSACTION;
-- Makes it easier to delete slots
ALTER TABLE slot_resources DROP CONSTRAINT slot_resources_slot_id_fkey;

ALTER TABLE slot_resources ALTER CONSTRAINT slot_resources_slot_id_fkey FOREIGN KEY (slot_id) REFERENCES slots(id) ON DELETE CASCADE;

UPDATE migrations
SET succeeded = true
WHERE migration = '011_slot_resources_FK_fix';
COMMIT;