-- The game value is used to differentiate between the mainline games. As the spin-offs (vita, psp)
-- derive from the same net code it seems to be required to specify, at least in the case of vita
ALTER TABLE slots ADD COLUMN IF NOT EXISTS game int NOT NULL DEFAULT 0;