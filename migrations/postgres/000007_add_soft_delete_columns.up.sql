-- Add soft delete columns to users table
ALTER TABLE users 
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN deleted_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- Add indexes for deleted_at
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;

-- Add is_active and soft delete columns to groups table
ALTER TABLE groups 
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN deleted_by UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_groups_deleted_at ON groups(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_groups_is_active ON groups(is_active);

-- Add is_active and soft delete columns to roles table
ALTER TABLE roles 
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN deleted_by UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_roles_deleted_at ON roles(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_is_active ON roles(is_active);

-- Add is_active and soft delete columns to permissions table
ALTER TABLE permissions 
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN deleted_by UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_permissions_deleted_at ON permissions(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_permissions_is_active ON permissions(is_active);

-- COMMENT: Partial indexes فقط record های active را index می‌کنند (performance)
