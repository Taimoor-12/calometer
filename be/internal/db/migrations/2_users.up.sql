BEGIN;

ALTER TABLE users DROP COLUMN email;

ALTER TABLE users ADD COLUMN username TEXT NOT NULL;

END;