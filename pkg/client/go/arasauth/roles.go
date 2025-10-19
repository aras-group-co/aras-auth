package arasauth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// CreateRoleRequest represents the request to create a role
type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateRoleRequest represents the request to update a role
type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListRolesResponse represents the response from listing roles
type ListRolesResponse struct {
	Roles []*Role `json:"roles"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
}

// AssignPermissionRequest represents the request to assign permission to role
type AssignPermissionRequest struct {
	PermissionID string `json:"permission_id"`
}

// CreateRole creates a new role
func (c *Client) CreateRole(ctx context.Context, name, description string) (*Role, error) {
	req := CreateRoleRequest{
		Name:        name,
		Description: description,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/roles", req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract role from data
	roleData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	role := &Role{}
	if id, ok := roleData["id"].(string); ok {
		role.ID = id
	}
	if name, ok := roleData["name"].(string); ok {
		role.Name = name
	}
	if description, ok := roleData["description"].(string); ok {
		role.Description = description
	}
	if createdAt, ok := roleData["created_at"].(string); ok {
		role.CreatedAt = createdAt
	}
	if updatedAt, ok := roleData["updated_at"].(string); ok {
		role.UpdatedAt = updatedAt
	}

	return role, nil
}

// ListRoles retrieves a list of roles with pagination
func (c *Client) ListRoles(ctx context.Context, page, limit int) (*ListRolesResponse, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	endpoint := "/api/v1/roles"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract roles response from data
	rolesData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	response := &ListRolesResponse{}
	if total, ok := rolesData["total"].(float64); ok {
		response.Total = int(total)
	}
	if page, ok := rolesData["page"].(float64); ok {
		response.Page = int(page)
	}
	if limit, ok := rolesData["limit"].(float64); ok {
		response.Limit = int(limit)
	}
	if roles, ok := rolesData["roles"].([]interface{}); ok {
		response.Roles = make([]*Role, len(roles))
		for i, roleData := range roles {
			if roleMap, ok := roleData.(map[string]interface{}); ok {
				role := &Role{}
				if id, ok := roleMap["id"].(string); ok {
					role.ID = id
				}
				if name, ok := roleMap["name"].(string); ok {
					role.Name = name
				}
				if description, ok := roleMap["description"].(string); ok {
					role.Description = description
				}
				if createdAt, ok := roleMap["created_at"].(string); ok {
					role.CreatedAt = createdAt
				}
				if updatedAt, ok := roleMap["updated_at"].(string); ok {
					role.UpdatedAt = updatedAt
				}
				response.Roles[i] = role
			}
		}
	}

	return response, nil
}

// GetRole retrieves a specific role by ID
func (c *Client) GetRole(ctx context.Context, roleID string) (*Role, error) {
	endpoint := fmt.Sprintf("/api/v1/roles/%s", roleID)

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract role from data
	roleData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	role := &Role{}
	if id, ok := roleData["id"].(string); ok {
		role.ID = id
	}
	if name, ok := roleData["name"].(string); ok {
		role.Name = name
	}
	if description, ok := roleData["description"].(string); ok {
		role.Description = description
	}
	if createdAt, ok := roleData["created_at"].(string); ok {
		role.CreatedAt = createdAt
	}
	if updatedAt, ok := roleData["updated_at"].(string); ok {
		role.UpdatedAt = updatedAt
	}

	return role, nil
}

// UpdateRole updates a role's information
func (c *Client) UpdateRole(ctx context.Context, roleID string, req *UpdateRoleRequest) (*Role, error) {
	endpoint := fmt.Sprintf("/api/v1/roles/%s", roleID)

	resp, err := c.makeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract role from data
	roleData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	role := &Role{}
	if id, ok := roleData["id"].(string); ok {
		role.ID = id
	}
	if name, ok := roleData["name"].(string); ok {
		role.Name = name
	}
	if description, ok := roleData["description"].(string); ok {
		role.Description = description
	}
	if createdAt, ok := roleData["created_at"].(string); ok {
		role.CreatedAt = createdAt
	}
	if updatedAt, ok := roleData["updated_at"].(string); ok {
		role.UpdatedAt = updatedAt
	}

	return role, nil
}

// DeleteRole deletes a role
func (c *Client) DeleteRole(ctx context.Context, roleID string) error {
	endpoint := fmt.Sprintf("/api/v1/roles/%s", roleID)

	resp, err := c.makeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return err
	}

	return nil
}

// AssignPermissionToRole assigns a permission to a role
func (c *Client) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	req := AssignPermissionRequest{
		PermissionID: permissionID,
	}

	endpoint := fmt.Sprintf("/api/v1/roles/%s/permissions", roleID)

	resp, err := c.makeRequest(ctx, "POST", endpoint, req)
	if err != nil {
		return err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return err
	}

	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (c *Client) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	endpoint := fmt.Sprintf("/api/v1/roles/%s/permissions/%s", roleID, permissionID)

	resp, err := c.makeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return err
	}

	return nil
}

// GetRolePermissions retrieves all permissions of a role
func (c *Client) GetRolePermissions(ctx context.Context, roleID string) ([]*Permission, error) {
	endpoint := fmt.Sprintf("/api/v1/roles/%s/permissions", roleID)

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract permissions from data
	permissionsData, ok := apiResp.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	permissions := make([]*Permission, len(permissionsData))
	for i, permissionData := range permissionsData {
		if permissionMap, ok := permissionData.(map[string]interface{}); ok {
			permission := &Permission{}
			if id, ok := permissionMap["id"].(string); ok {
				permission.ID = id
			}
			if resource, ok := permissionMap["resource"].(string); ok {
				permission.Resource = resource
			}
			if action, ok := permissionMap["action"].(string); ok {
				permission.Action = action
			}
			if description, ok := permissionMap["description"].(string); ok {
				permission.Description = description
			}
			if createdAt, ok := permissionMap["created_at"].(string); ok {
				permission.CreatedAt = createdAt
			}
			if updatedAt, ok := permissionMap["updated_at"].(string); ok {
				permission.UpdatedAt = updatedAt
			}
			permissions[i] = permission
		}
	}

	return permissions, nil
}
