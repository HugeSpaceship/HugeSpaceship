ALTER TABLE users ADD COLUMN IF NOT EXISTS planet_lbp2 varchar(40) default '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS planet_lbp3 varchar(40) default '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS planet_lbp_vita varchar(40) default '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS planet_cc varchar(40) default '';

ALTER TABLE users ADD COLUMN IF NOT EXISTS boo_icon varchar(40) default '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS meh_icon varchar(40) default '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS yay_icon varchar(40) default '';

ALTER TABLE users RENAME locationx TO location_x;
ALTER TABLE users RENAME locationy TO location_y;

ALTER TABLE slots RENAME locationx TO location_x;
ALTER TABLE slots RENAME locationy TO location_y;