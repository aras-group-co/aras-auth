package arasauth

import (
	"context"
	"fmt"
)

// Login authenticates a user and returns tokens
func (c *Client) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/auth/login", req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract auth response from data
	authData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	authResp := &AuthResponse{}
	if accessToken, ok := authData["access_token"].(string); ok {
		authResp.AccessToken = accessToken
	}
	if refreshToken, ok := authData["refresh_token"].(string); ok {
		authResp.RefreshToken = refreshToken
	}
	if expiresIn, ok := authData["expires_in"].(float64); ok {
		authResp.ExpiresIn = int64(expiresIn)
	}
	if tokenType, ok := authData["token_type"].(string); ok {
		authResp.TokenType = tokenType
	}
	if userData, ok := authData["user"].(map[string]interface{}); ok {
		authResp.User = &User{}
		if id, ok := userData["id"].(string); ok {
			authResp.User.ID = id
		}
		if email, ok := userData["email"].(string); ok {
			authResp.User.Email = email
		}
		if firstName, ok := userData["first_name"].(string); ok {
			authResp.User.FirstName = firstName
		}
		if lastName, ok := userData["last_name"].(string); ok {
			authResp.User.LastName = lastName
		}
		if status, ok := userData["status"].(string); ok {
			authResp.User.Status = status
		}
		if emailVerified, ok := userData["email_verified"].(bool); ok {
			authResp.User.EmailVerified = emailVerified
		}
		if createdAt, ok := userData["created_at"].(string); ok {
			authResp.User.CreatedAt = createdAt
		}
		if updatedAt, ok := userData["updated_at"].(string); ok {
			authResp.User.UpdatedAt = updatedAt
		}
	}

	// Set token for future requests
	c.SetToken(authResp.AccessToken)

	return authResp, nil
}

// Register creates a new user account
func (c *Client) Register(ctx context.Context, email, password, firstName, lastName string) (*User, error) {
	req := RegisterRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/auth/register", req)
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

// RefreshToken refreshes the access token using a refresh token
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	req := map[string]string{
		"refresh_token": refreshToken,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/auth/refresh", req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract auth response from data
	authData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	authResp := &AuthResponse{}
	if accessToken, ok := authData["access_token"].(string); ok {
		authResp.AccessToken = accessToken
	}
	if refreshToken, ok := authData["refresh_token"].(string); ok {
		authResp.RefreshToken = refreshToken
	}
	if expiresIn, ok := authData["expires_in"].(float64); ok {
		authResp.ExpiresIn = int64(expiresIn)
	}
	if tokenType, ok := authData["token_type"].(string); ok {
		authResp.TokenType = tokenType
	}
	if userData, ok := authData["user"].(map[string]interface{}); ok {
		authResp.User = &User{}
		if id, ok := userData["id"].(string); ok {
			authResp.User.ID = id
		}
		if email, ok := userData["email"].(string); ok {
			authResp.User.Email = email
		}
		if firstName, ok := userData["first_name"].(string); ok {
			authResp.User.FirstName = firstName
		}
		if lastName, ok := userData["last_name"].(string); ok {
			authResp.User.LastName = lastName
		}
		if status, ok := userData["status"].(string); ok {
			authResp.User.Status = status
		}
		if emailVerified, ok := userData["email_verified"].(bool); ok {
			authResp.User.EmailVerified = emailVerified
		}
		if createdAt, ok := userData["created_at"].(string); ok {
			authResp.User.CreatedAt = createdAt
		}
		if updatedAt, ok := userData["updated_at"].(string); ok {
			authResp.User.UpdatedAt = updatedAt
		}
	}

	// Set token for future requests
	c.SetToken(authResp.AccessToken)

	return authResp, nil
}

// Logout logs out the user by invalidating the refresh token
func (c *Client) Logout(ctx context.Context, refreshToken string) error {
	req := map[string]string{
		"refresh_token": refreshToken,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/auth/logout", req)
	if err != nil {
		return err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return err
	}

	// Clear token
	c.SetToken("")

	return nil
}

// GetCurrentUser gets the current authenticated user
func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	resp, err := c.makeRequest(ctx, "GET", "/api/v1/users/me", nil)
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

// CheckPermission checks if a user has a specific permission
func (c *Client) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	req := CheckPermissionRequest{
		UserID:   userID,
		Resource: resource,
		Action:   action,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/authz/check", req)
	if err != nil {
		return false, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return false, err
	}

	// Extract permission check result from data
	permissionData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("invalid response format")
	}

	hasPermission, ok := permissionData["has_permission"].(bool)
	if !ok {
		return false, fmt.Errorf("invalid response format")
	}

	return hasPermission, nil
}

// ChangePassword changes the current user's password
func (c *Client) ChangePassword(ctx context.Context, currentPassword, newPassword string) error {
	req := ChangePasswordRequest{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/auth/change-password", req)
	if err != nil {
		return err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return err
	}

	return nil
}

// ForgotPassword requests a password reset email
func (c *Client) ForgotPassword(ctx context.Context, email string) error {
	req := ForgotPasswordRequest{
		Email: email,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/auth/forgot-password", req)
	if err != nil {
		return err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return err
	}

	return nil
}

// ResetPassword resets password using a reset token
func (c *Client) ResetPassword(ctx context.Context, token, newPassword string) error {
	req := ResetPasswordRequest{
		Token:       token,
		NewPassword: newPassword,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/auth/reset-password", req)
	if err != nil {
		return err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return err
	}

	return nil
}

// VerifyEmail verifies a user's email address
func (c *Client) VerifyEmail(ctx context.Context, userID string) error {
	req := VerifyEmailRequest{
		UserID: userID,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/auth/verify-email", req)
	if err != nil {
		return err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return err
	}

	return nil
}

// IntrospectToken introspects a token and returns its information
func (c *Client) IntrospectToken(ctx context.Context, token string) (*TokenIntrospection, error) {
	req := map[string]string{
		"token": token,
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/auth/introspect", req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	// Extract token introspection from data
	introspectionData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	introspection := &TokenIntrospection{}
	if active, ok := introspectionData["active"].(bool); ok {
		introspection.Active = active
	}
	if userID, ok := introspectionData["user_id"].(string); ok {
		introspection.UserID = userID
	}
	if email, ok := introspectionData["email"].(string); ok {
		introspection.Email = email
	}
	if expiresAt, ok := introspectionData["exp"].(float64); ok {
		introspection.ExpiresAt = int64(expiresAt)
	}

	return introspection, nil
}
