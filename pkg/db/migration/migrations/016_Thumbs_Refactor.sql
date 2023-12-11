-- Consolidate thumbs tables.

ALTER TABLE thumbs_up RENAME TO thumbs;
DROP TABLE thumbs_down;
ALTER TABLE thumbs ADD COLUMN down bool;