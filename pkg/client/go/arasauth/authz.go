package arasauth

import (
	"context"
	"fmt"
)

// AssignRoleRequest represents the request to assign a role
type AssignRoleRequest struct {
	RoleID string `json:"role_id"`
}

// AssignRoleToUser assigns a role to a user
func (c *Client) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	req := AssignRoleRequest{
		RoleID: roleID,
	}

	endpoint := fmt.Sprintf("/api/v1/users/%s/roles", userID)

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

// RemoveRoleFromUser removes a role from a user
func (c *Client) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	endpoint := fmt.Sprintf("/api/v1/users/%s/roles/%s", userID, roleID)

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

// GetUserRoles retrieves all roles of a user
func (c *Client) GetUserRoles(ctx context.Context, userID string) ([]*Role, error) {
	endpoint := fmt.Sprintf("/api/v1/users/%s/roles", userID)

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract roles from data
	rolesData, ok := apiResp.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	roles := make([]*Role, len(rolesData))
	for i, roleData := range rolesData {
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
			roles[i] = role
		}
	}

	return roles, nil
}

// AssignRoleToGroup assigns a role to a group
func (c *Client) AssignRoleToGroup(ctx context.Context, groupID, roleID string) error {
	req := AssignRoleRequest{
		RoleID: roleID,
	}

	endpoint := fmt.Sprintf("/api/v1/groups/%s/roles", groupID)

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

// RemoveRoleFromGroup removes a role from a group
func (c *Client) RemoveRoleFromGroup(ctx context.Context, groupID, roleID string) error {
	endpoint := fmt.Sprintf("/api/v1/groups/%s/roles/%s", groupID, roleID)

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

// GetGroupRoles retrieves all roles of a group
func (c *Client) GetGroupRoles(ctx context.Context, groupID string) ([]*Role, error) {
	endpoint := fmt.Sprintf("/api/v1/groups/%s/roles", groupID)

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract roles from data
	rolesData, ok := apiResp.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	roles := make([]*Role, len(rolesData))
	for i, roleData := range rolesData {
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
			roles[i] = role
		}
	}

	return roles, nil
}
