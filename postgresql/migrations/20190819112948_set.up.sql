ALTER TABLE users ALTER COLUMN bdate SET DEFAULT '0001-01-01 00:00:00';
UPDATE users SET bdate='0001-01-01 00:00:00' WHERE bdate IS NULL;