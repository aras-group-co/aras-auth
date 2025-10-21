-- Add is_deleted to users table
ALTER TABLE users 
    ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

-- Add index for is_deleted
CREATE INDEX idx_users_is_deleted ON users(is_deleted) WHERE is_deleted = FALSE;

-- Add is_active and is_deleted to groups table
ALTER TABLE groups 
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_groups_is_deleted ON groups(is_deleted) WHERE is_deleted = FALSE;
CREATE INDEX idx_groups_is_active ON groups(is_active);

-- Add is_active and is_deleted to roles table
ALTER TABLE roles 
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_roles_is_deleted ON roles(is_deleted) WHERE is_deleted = FALSE;
CREATE INDEX idx_roles_is_active ON roles(is_active);

-- Add is_active and is_deleted to permissions table
ALTER TABLE permissions 
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_permissions_is_deleted ON permissions(is_deleted) WHERE is_deleted = FALSE;
CREATE INDEX idx_permissions_is_active ON permissions(is_active);

-- COMMENT: Partial indexes فقط record های active را index می‌کنند (performance)
