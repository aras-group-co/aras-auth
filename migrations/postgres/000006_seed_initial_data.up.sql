-- Insert default local provider
INSERT INTO providers (name, type, config, enabled) VALUES 
('local', 'local', '{}', TRUE);

-- Insert basic roles
INSERT INTO roles (name, description) VALUES 
('admin', 'System administrator with full access'),
('user', 'Regular user with basic access'),
('moderator', 'Moderator with elevated permissions');

-- Insert basic permissions
INSERT INTO permissions (resource, action, description) VALUES 
('users', 'create', 'Create new users'),
('users', 'read', 'View user information'),
('users', 'update', 'Update user information'),
('users', 'delete', 'Delete users'),
('groups', 'create', 'Create new groups'),
('groups', 'read', 'View group information'),
('groups', 'update', 'Update group information'),
('groups', 'delete', 'Delete groups'),
('roles', 'create', 'Create new roles'),
('roles', 'read', 'View role information'),
('roles', 'update', 'Update role information'),
('roles', 'delete', 'Delete roles'),
('permissions', 'create', 'Create new permissions'),
('permissions', 'read', 'View permission information'),
('permissions', 'update', 'Update permission information'),
('permissions', 'delete', 'Delete permissions');

-- Assign permissions to admin role (all permissions)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin';

-- Assign permissions to user role (basic read permissions)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'user' 
AND p.resource = 'users' AND p.action = 'read';

-- Assign permissions to moderator role (user and group management)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'moderator' 
AND ((p.resource = 'users' AND p.action IN ('read', 'update')) 
     OR (p.resource = 'groups' AND p.action IN ('read', 'update')));

