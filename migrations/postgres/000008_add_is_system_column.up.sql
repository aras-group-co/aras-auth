-- Add is_system column to users table
ALTER TABLE users 
    ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

-- Add index for is_system (optional, for performance)
CREATE INDEX idx_users_is_system ON users(is_system);

-- Add is_system column to providers table
ALTER TABLE providers 
    ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_providers_is_system ON providers(is_system);

-- Add is_system column to roles table
ALTER TABLE roles 
    ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_roles_is_system ON roles(is_system);

-- Add is_system column to permissions table
ALTER TABLE permissions 
    ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_permissions_is_system ON permissions(is_system);

-- Add is_system column to groups table
ALTER TABLE groups 
    ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_groups_is_system ON groups(is_system);

-- Update existing seed data to is_system = TRUE
-- Provider
UPDATE providers SET is_system = TRUE WHERE name = 'local';

-- Roles
UPDATE roles SET is_system = TRUE WHERE name IN ('admin', 'user', 'moderator');

-- Permissions (all 16 seed permissions)
UPDATE permissions SET is_system = TRUE WHERE resource IN ('users', 'groups', 'roles', 'permissions');

-- Admin user (if exists from previous migration)
UPDATE users SET is_system = TRUE WHERE email = 'admin@aras-services.com';
