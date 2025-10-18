package arasauth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListGroupsResponse represents the response from listing groups
type ListGroupsResponse struct {
	Groups []*Group `json:"groups"`
	Total  int      `json:"total"`
	Page   int      `json:"page"`
	Limit  int      `json:"limit"`
}

// CreateGroupRequest represents the request to create a group
type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateGroupRequest represents the request to update a group
type UpdateGroupRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// AddMemberRequest represents the request to add a member to a group
type AddMemberRequest struct {
	UserID string `json:"user_id"`
}

// ListGroups retrieves a list of groups with pagination
func (c *Client) ListGroups(ctx context.Context, page, limit int) (*ListGroupsResponse, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	endpoint := "/api/v1/groups"
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

	// Extract groups response from data
	groupsData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	response := &ListGroupsResponse{}
	if total, ok := groupsData["total"].(float64); ok {
		response.Total = int(total)
	}
	if page, ok := groupsData["page"].(float64); ok {
		response.Page = int(page)
	}
	if limit, ok := groupsData["limit"].(float64); ok {
		response.Limit = int(limit)
	}
	if groups, ok := groupsData["groups"].([]interface{}); ok {
		response.Groups = make([]*Group, len(groups))
		for i, groupData := range groups {
			if groupMap, ok := groupData.(map[string]interface{}); ok {
				group := &Group{}
				if id, ok := groupMap["id"].(string); ok {
					group.ID = id
				}
				if name, ok := groupMap["name"].(string); ok {
					group.Name = name
				}
				if description, ok := groupMap["description"].(string); ok {
					group.Description = description
				}
				if createdAt, ok := groupMap["created_at"].(string); ok {
					group.CreatedAt = createdAt
				}
				if updatedAt, ok := groupMap["updated_at"].(string); ok {
					group.UpdatedAt = updatedAt
				}
				response.Groups[i] = group
			}
		}
	}

	return response, nil
}

// CreateGroup creates a new group
func (c *Client) CreateGroup(ctx context.Context, name, description string) (*Group, error) {
	req := CreateGroupRequest{
		Name:        name,
		Description: description,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/groups", req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract group from data
	groupData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	group := &Group{}
	if id, ok := groupData["id"].(string); ok {
		group.ID = id
	}
	if name, ok := groupData["name"].(string); ok {
		group.Name = name
	}
	if description, ok := groupData["description"].(string); ok {
		group.Description = description
	}
	if createdAt, ok := groupData["created_at"].(string); ok {
		group.CreatedAt = createdAt
	}
	if updatedAt, ok := groupData["updated_at"].(string); ok {
		group.UpdatedAt = updatedAt
	}

	return group, nil
}

// GetGroup retrieves a specific group by ID
func (c *Client) GetGroup(ctx context.Context, groupID string) (*Group, error) {
	endpoint := fmt.Sprintf("/api/v1/groups/%s", groupID)

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract group from data
	groupData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	group := &Group{}
	if id, ok := groupData["id"].(string); ok {
		group.ID = id
	}
	if name, ok := groupData["name"].(string); ok {
		group.Name = name
	}
	if description, ok := groupData["description"].(string); ok {
		group.Description = description
	}
	if createdAt, ok := groupData["created_at"].(string); ok {
		group.CreatedAt = createdAt
	}
	if updatedAt, ok := groupData["updated_at"].(string); ok {
		group.UpdatedAt = updatedAt
	}

	return group, nil
}

// UpdateGroup updates a group's information
func (c *Client) UpdateGroup(ctx context.Context, groupID string, req *UpdateGroupRequest) (*Group, error) {
	endpoint := fmt.Sprintf("/api/v1/groups/%s", groupID)

	resp, err := c.makeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract group from data
	groupData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	group := &Group{}
	if id, ok := groupData["id"].(string); ok {
		group.ID = id
	}
	if name, ok := groupData["name"].(string); ok {
		group.Name = name
	}
	if description, ok := groupData["description"].(string); ok {
		group.Description = description
	}
	if createdAt, ok := groupData["created_at"].(string); ok {
		group.CreatedAt = createdAt
	}
	if updatedAt, ok := groupData["updated_at"].(string); ok {
		group.UpdatedAt = updatedAt
	}

	return group, nil
}

// DeleteGroup deletes a group
func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	endpoint := fmt.Sprintf("/api/v1/groups/%s", groupID)

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

// AddMember adds a user to a group
func (c *Client) AddMember(ctx context.Context, groupID, userID string) error {
	req := AddMemberRequest{
		UserID: userID,
	}

	endpoint := fmt.Sprintf("/api/v1/groups/%s/members", groupID)

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

// RemoveMember removes a user from a group
func (c *Client) RemoveMember(ctx context.Context, groupID, userID string) error {
	endpoint := fmt.Sprintf("/api/v1/groups/%s/members/%s", groupID, userID)

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

// GetMembers retrieves all members of a group
func (c *Client) GetMembers(ctx context.Context, groupID string) ([]*User, error) {
	endpoint := fmt.Sprintf("/api/v1/groups/%s/members", groupID)

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract users from data
	usersData, ok := apiResp.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	users := make([]*User, len(usersData))
	for i, userData := range usersData {
		if userMap, ok := userData.(map[string]interface{}); ok {
			user := &User{}
			if id, ok := userMap["id"].(string); ok {
				user.ID = id
			}
			if email, ok := userMap["email"].(string); ok {
				user.Email = email
			}
			if firstName, ok := userMap["first_name"].(string); ok {
				user.FirstName = firstName
			}
			if lastName, ok := userMap["last_name"].(string); ok {
				user.LastName = lastName
			}
			if status, ok := userMap["status"].(string); ok {
				user.Status = status
			}
			if emailVerified, ok := userMap["email_verified"].(bool); ok {
				user.EmailVerified = emailVerified
			}
			if createdAt, ok := userMap["created_at"].(string); ok {
				user.CreatedAt = createdAt
			}
			if updatedAt, ok := userMap["updated_at"].(string); ok {
				user.UpdatedAt = updatedAt
			}
			users[i] = user
		}
	}

	return users, nil
}


