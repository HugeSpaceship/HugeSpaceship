-- Makes it easier to delete slots
ALTER TABLE slot_resources DROP CONSTRAINT slot_resources_slot_id_fkey;

ALTER TABLE slot_resources ADD CONSTRAINT slot_resources_slot_id_fkey  FOREIGN KEY (slot_id) REFERENCES slots(id) ON DELETE CASCADE;