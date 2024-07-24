CREATE TYPE resource_override AS ENUM (
    'ALLOW',
    'DENY',
    'EXTERNAL'
    );

CREATE TABLE resource_overrides (
    hash text not null,
    override resource_override not null,
    reason text
);