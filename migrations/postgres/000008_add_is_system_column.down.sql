-- Remove is_system column from all tables
ALTER TABLE users DROP COLUMN IF EXISTS is_system;
ALTER TABLE providers DROP COLUMN IF EXISTS is_system;
ALTER TABLE roles DROP COLUMN IF EXISTS is_system;
ALTER TABLE permissions DROP COLUMN IF EXISTS is_system;
ALTER TABLE groups DROP COLUMN IF EXISTS is_system;
