package arasauth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListUsersResponse represents the response from listing users
type ListUsersResponse struct {
	Users []*User `json:"users"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Status    *string `json:"status,omitempty"`
}

// ListUsers retrieves a list of users with pagination
func (c *Client) ListUsers(ctx context.Context, page, limit int) (*ListUsersResponse, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	endpoint := "/api/v1/users"
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

	// Extract users response from data
	usersData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	response := &ListUsersResponse{}
	if total, ok := usersData["total"].(float64); ok {
		response.Total = int(total)
	}
	if page, ok := usersData["page"].(float64); ok {
		response.Page = int(page)
	}
	if limit, ok := usersData["limit"].(float64); ok {
		response.Limit = int(limit)
	}
	if users, ok := usersData["users"].([]interface{}); ok {
		response.Users = make([]*User, len(users))
		for i, userData := range users {
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
				response.Users[i] = user
			}
		}
	}

	return response, nil
}

// GetUser retrieves a specific user by ID
func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	endpoint := fmt.Sprintf("/api/v1/users/%s", userID)

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract user from data
	userData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	user := &User{}
	if id, ok := userData["id"].(string); ok {
		user.ID = id
	}
	if email, ok := userData["email"].(string); ok {
		user.Email = email
	}
	if firstName, ok := userData["first_name"].(string); ok {
		user.FirstName = firstName
	}
	if lastName, ok := userData["last_name"].(string); ok {
		user.LastName = lastName
	}
	if status, ok := userData["status"].(string); ok {
		user.Status = status
	}
	if emailVerified, ok := userData["email_verified"].(bool); ok {
		user.EmailVerified = emailVerified
	}
	if createdAt, ok := userData["created_at"].(string); ok {
		user.CreatedAt = createdAt
	}
	if updatedAt, ok := userData["updated_at"].(string); ok {
		user.UpdatedAt = updatedAt
	}

	return user, nil
}

// UpdateUser updates a user's information
func (c *Client) UpdateUser(ctx context.Context, userID string, req *UpdateUserRequest) (*User, error) {
	endpoint := fmt.Sprintf("/api/v1/users/%s", userID)

	resp, err := c.makeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract user from data
	userData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	user := &User{}
	if id, ok := userData["id"].(string); ok {
		user.ID = id
	}
	if email, ok := userData["email"].(string); ok {
		user.Email = email
	}
	if firstName, ok := userData["first_name"].(string); ok {
		user.FirstName = firstName
	}
	if lastName, ok := userData["last_name"].(string); ok {
		user.LastName = lastName
	}
	if status, ok := userData["status"].(string); ok {
		user.Status = status
	}
	if emailVerified, ok := userData["email_verified"].(bool); ok {
		user.EmailVerified = emailVerified
	}
	if createdAt, ok := userData["created_at"].(string); ok {
		user.CreatedAt = createdAt
	}
	if updatedAt, ok := userData["updated_at"].(string); ok {
		user.UpdatedAt = updatedAt
	}

	return user, nil
}

// DeleteUser deletes a user
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	endpoint := fmt.Sprintf("/api/v1/users/%s", userID)

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


