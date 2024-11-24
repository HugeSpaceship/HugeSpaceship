CREATE TYPE override_type AS ENUM (
    'ALLOW',
    'DENY',
    'EXTERNAL'
    );

CREATE TABLE resource_overrides (
    hash text not null,
    override override_type not null,
    reason text
);