BEGIN TRANSACTION;

INSERT INTO migrations (migration, succeeded)
VALUES ('002_Session_Changes', false);

ALTER TABLE sessions DROP COLUMN id;
DROP SEQUENCE sessions_id_seq;
ALTER TABLE sessions ADD PRIMARY KEY (token);
ALTER TABLE sessions ADD CONSTRAINT Unique_Token UNIQUE (token);

CREATE TABLE resources (
    Hash char[40],
    File bytea, -- blob
    OriginalUploader int references users(id)
);

UPDATE migrations
SET succeeded = true
WHERE migration = '002_Session_Changes';

COMMIT;