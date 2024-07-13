CREATE TYPE privacy AS ENUM ( -- what platform are they on,
    'all',
    'psn',
    'game'
    );

ALTER TABLE users ADD COLUMN level_visibility privacy;
ALTER TABLE users ADD COLUMN profile_visibility privacy;