-- Add published column to facilitate re-uploading of levels.
ALTER TABLE slots ADD COLUMN published bool DEFAULT true;