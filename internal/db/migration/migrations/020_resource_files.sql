CREATE TABLE files (
    hash text PRIMARY KEY,
    file oid NOT NULL
);

ALTER TABLE resources RENAME TO resources_old;

CREATE TYPE resource_backends AS ENUM (
    'file',
    'pg_lob',
    's3'
    );

CREATE TABLE resources (
    hash text PRIMARY KEY,
    backend resource_backends not null default 'pg_lob',
    resource_type resource_type NOT NULL DEFAULT 'UNK',
    size bigint,
    created timestamptz,
    uploader uuid references users(id)
);

INSERT INTO files(hash, file) SELECT hash, file FROM resources_old;

INSERT INTO resources(hash, resource_type, size, created, uploader)
    SELECT hash, resource_type, size, created, originaluploader FROM resources_old;

ALTER TABLE slot_resources DROP CONSTRAINT slot_resources_resource_hash_fkey;
ALTER TABLE slot_resources ADD CONSTRAINT slot_resources_resource_hash_fkey FOREIGN KEY(resource_hash) REFERENCES resources(hash);

ALTER TABLE slots DROP CONSTRAINT slots_root_level_fkey;
ALTER TABLE slots ADD CONSTRAINT slots_root_level_fkey FOREIGN KEY(root_level) REFERENCES resources(hash);


DROP TABLE resources_old;
