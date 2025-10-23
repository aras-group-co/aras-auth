package arasauth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// CreatePermissionRequest represents the request to create a permission
type CreatePermissionRequest struct {
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// UpdatePermissionRequest represents the request to update a permission
type UpdatePermissionRequest struct {
	Resource    *string `json:"resource,omitempty"`
	Action      *string `json:"action,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// ListPermissionsResponse represents the response from listing permissions
type ListPermissionsResponse struct {
	Permissions []*Permission `json:"permissions"`
	Total       int           `json:"total"`
	Page        int           `json:"page"`
	Limit       int           `json:"limit"`
}

// CreatePermission creates a new permission
func (c *Client) CreatePermission(ctx context.Context, resource, action, description string, isActive *bool) (*Permission, error) {
	req := CreatePermissionRequest{
		Resource:    resource,
		Action:      action,
		Description: description,
		IsActive:    isActive,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/permissions", req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract permission from data
	permissionData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	permission := &Permission{}
	if id, ok := permissionData["id"].(string); ok {
		permission.ID = id
	}
	if resource, ok := permissionData["resource"].(string); ok {
		permission.Resource = resource
	}
	if action, ok := permissionData["action"].(string); ok {
		permission.Action = action
	}
	if description, ok := permissionData["description"].(string); ok {
		permission.Description = description
	}
	if isActive, ok := permissionData["is_active"].(bool); ok {
		permission.IsActive = isActive
	}
	if isDeleted, ok := permissionData["is_deleted"].(bool); ok {
		permission.IsDeleted = isDeleted
	}
	if isSystem, ok := permissionData["is_system"].(bool); ok {
		permission.IsSystem = isSystem
	}
	if createdAt, ok := permissionData["created_at"].(string); ok {
		permission.CreatedAt = createdAt
	}
	if updatedAt, ok := permissionData["updated_at"].(string); ok {
		permission.UpdatedAt = updatedAt
	}

	return permission, nil
}

// ListPermissions retrieves a list of permissions with pagination
func (c *Client) ListPermissions(ctx context.Context, page, limit int) (*ListPermissionsResponse, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	endpoint := "/api/v1/permissions"
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

	// Extract permissions response from data
	permissionsData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	response := &ListPermissionsResponse{}
	if total, ok := permissionsData["total"].(float64); ok {
		response.Total = int(total)
	}
	if page, ok := permissionsData["page"].(float64); ok {
		response.Page = int(page)
	}
	if limit, ok := permissionsData["limit"].(float64); ok {
		response.Limit = int(limit)
	}
	if permissions, ok := permissionsData["permissions"].([]interface{}); ok {
		response.Permissions = make([]*Permission, len(permissions))
		for i, permissionData := range permissions {
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
				response.Permissions[i] = permission
			}
		}
	}

	return response, nil
}

// GetPermission retrieves a specific permission by ID
func (c *Client) GetPermission(ctx context.Context, permissionID string) (*Permission, error) {
	endpoint := fmt.Sprintf("/api/v1/permissions/%s", permissionID)

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract permission from data
	permissionData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	permission := &Permission{}
	if id, ok := permissionData["id"].(string); ok {
		permission.ID = id
	}
	if resource, ok := permissionData["resource"].(string); ok {
		permission.Resource = resource
	}
	if action, ok := permissionData["action"].(string); ok {
		permission.Action = action
	}
	if description, ok := permissionData["description"].(string); ok {
		permission.Description = description
	}
	if isActive, ok := permissionData["is_active"].(bool); ok {
		permission.IsActive = isActive
	}
	if isDeleted, ok := permissionData["is_deleted"].(bool); ok {
		permission.IsDeleted = isDeleted
	}
	if isSystem, ok := permissionData["is_system"].(bool); ok {
		permission.IsSystem = isSystem
	}
	if createdAt, ok := permissionData["created_at"].(string); ok {
		permission.CreatedAt = createdAt
	}
	if updatedAt, ok := permissionData["updated_at"].(string); ok {
		permission.UpdatedAt = updatedAt
	}

	return permission, nil
}

// UpdatePermission updates a permission's information
func (c *Client) UpdatePermission(ctx context.Context, permissionID string, req *UpdatePermissionRequest) (*Permission, error) {
	endpoint := fmt.Sprintf("/api/v1/permissions/%s", permissionID)

	resp, err := c.makeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract permission from data
	permissionData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	permission := &Permission{}
	if id, ok := permissionData["id"].(string); ok {
		permission.ID = id
	}
	if resource, ok := permissionData["resource"].(string); ok {
		permission.Resource = resource
	}
	if action, ok := permissionData["action"].(string); ok {
		permission.Action = action
	}
	if description, ok := permissionData["description"].(string); ok {
		permission.Description = description
	}
	if isActive, ok := permissionData["is_active"].(bool); ok {
		permission.IsActive = isActive
	}
	if isDeleted, ok := permissionData["is_deleted"].(bool); ok {
		permission.IsDeleted = isDeleted
	}
	if isSystem, ok := permissionData["is_system"].(bool); ok {
		permission.IsSystem = isSystem
	}
	if createdAt, ok := permissionData["created_at"].(string); ok {
		permission.CreatedAt = createdAt
	}
	if updatedAt, ok := permissionData["updated_at"].(string); ok {
		permission.UpdatedAt = updatedAt
	}

	return permission, nil
}

// DeletePermission deletes a permission
func (c *Client) DeletePermission(ctx context.Context, permissionID string) error {
	endpoint := fmt.Sprintf("/api/v1/permissions/%s", permissionID)

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
