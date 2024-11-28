-- Turns out that resources most likely doesn't need to be this complicated

ALTER TABLE resources DROP COLUMN backend;
ALTER TABLE resources DROP COLUMN backend_name;