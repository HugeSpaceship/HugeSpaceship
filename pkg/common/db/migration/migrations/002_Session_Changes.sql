BEGIN TRANSACTION;

INSERT INTO migrations (migration, succeeded)
VALUES ('002_Session_Changes', false);

ALTER TABLE sessions DROP COLUMN IF EXISTS id;
ALTER TABLE sessions ADD PRIMARY KEY (token);
ALTER TABLE sessions ADD CONSTRAINT Unique_Token UNIQUE (token);

CREATE TABLE resources (
    Hash char(40) PRIMARY KEY,
    Size bigint NOT NULL,
    File oid, -- blob
    OriginalUploader uuid references users(id)
);

UPDATE migrations
SET succeeded = true
WHERE migration = '002_Session_Changes';

COMMIT;