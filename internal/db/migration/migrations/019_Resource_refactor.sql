CREATE TYPE resource_type AS ENUM (
    'TEX', -- Texture
    'PLN', -- Plan
    'PNG', -- PNG image
    'JPG', -- Jpeg image
    'REC', -- Motion Recording
    'FSH', -- Fish Finger Script
    'VOP', -- Sound (voice)
    'LVL', -- Level
    'ADC', -- LBP3 Adventure
    'ADS', -- LBP3 Adventure
    'QST', -- Quest
    'CHK', -- Streaming chunk
    'UNK' -- Unknown
    );

ALTER TABLE resources ADD COLUMN resource_type resource_type NOT NULL DEFAULT 'UNK';
ALTER TABLE resources ADD COLUMN created timestamptz;