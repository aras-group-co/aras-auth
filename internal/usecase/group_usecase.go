package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/aras-services/aras-auth/internal/domain"
)

type GroupUseCase struct {
	groupRepo domain.GroupRepository
}

func NewGroupUseCase(groupRepo domain.GroupRepository) *GroupUseCase {
	return &GroupUseCase{
		groupRepo: groupRepo,
	}
}

type ListGroupsResponse struct {
	Groups []*domain.Group `json:"groups"`
	Total  int             `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

func (uc *GroupUseCase) CreateGroup(ctx context.Context, req *domain.CreateGroupRequest) (*domain.Group, error) {
	isActive := true // default
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	group := &domain.Group{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		IsActive:    isActive,
		IsDeleted:   false,
		IsSystem:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.groupRepo.Create(group); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return group, nil
}

func (uc *GroupUseCase) GetGroup(ctx context.Context, groupID uuid.UUID) (*domain.Group, error) {
	return uc.groupRepo.GetByID(groupID)
}

func (uc *GroupUseCase) ListGroups(ctx context.Context, page, limit int) (*ListGroupsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	groups, err := uc.groupRepo.List(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	total, err := uc.groupRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("failed to count groups: %w", err)
	}

	return &ListGroupsResponse{
		Groups: groups,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

func (uc *GroupUseCase) UpdateGroup(ctx context.Context, groupID uuid.UUID, req *domain.UpdateGroupRequest) (*domain.Group, error) {
	// Get existing group
	group, err := uc.groupRepo.GetByID(groupID)
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.Description != nil {
		group.Description = *req.Description
	}
	if req.IsActive != nil {
		group.IsActive = *req.IsActive
	}

	// Save updated group
	if err := uc.groupRepo.Update(group); err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	return group, nil
}

func (uc *GroupUseCase) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	// Check if group exists
	_, err := uc.groupRepo.GetByID(groupID)
	if err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	// Delete group
	return uc.groupRepo.Delete(groupID)
}

func (uc *GroupUseCase) AddMember(ctx context.Context, groupID uuid.UUID, req *domain.AddMemberRequest) error {
	return uc.groupRepo.AddMember(groupID, req.UserID)
}

func (uc *GroupUseCase) RemoveMember(ctx context.Context, groupID, userID uuid.UUID) error {
	return uc.groupRepo.RemoveMember(groupID, userID)
}

func (uc *GroupUseCase) GetMembers(ctx context.Context, groupID uuid.UUID) ([]*domain.User, error) {
	return uc.groupRepo.GetMembers(groupID)
}

func (uc *GroupUseCase) GetUserGroups(ctx context.Context, userID uuid.UUID) ([]*domain.Group, error) {
	return uc.groupRepo.GetUserGroups(userID)
}
