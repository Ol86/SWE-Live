package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"SWE-Live/internal/repository"
)

var (
	ErrMemberNotFound       = errors.New("member not found")
	ErrOptimisticLockFailed = errors.New("member was changed concurrently")
)

// MemberWriteService handles create, update, and delete operations for members.
type MemberWriteService interface {
	Create(ctx context.Context, cmd CreateMemberCommand) (MemberWriteModel, error)
	Update(ctx context.Context, cmd UpdateMemberCommand) (MemberWriteModel, error)
	Delete(ctx context.Context, id int32) error
}

// DefaultMemberWriteService is the default implementation of MemberWriteService.
type DefaultMemberWriteService struct {
	members repository.MemberRepository
}

// NewMemberWriteService creates a new member write service.
func NewMemberWriteService(members repository.MemberRepository) *DefaultMemberWriteService {
	return &DefaultMemberWriteService{members: members}
}

// CreateMemberCommand contains the values needed to create a new member.
type CreateMemberCommand struct {
	Username     string  `json:"username"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Gender       *string `json:"gender,omitempty"`
	DateOfBirth  string  `json:"dateOfBirth"`
	MemberSince  *string `json:"memberSince,omitempty"`
	IsStudent    *bool   `json:"isStudent,omitempty"`
	EmailAddress string  `json:"emailAddress"`
	Interests    []byte  `json:"interests,omitempty"`
}

// UpdateMemberCommand contains the values needed to update an existing member.
type UpdateMemberCommand struct {
	ID           int32   `json:"id"`
	Version      int32   `json:"version"`
	Username     string  `json:"username"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Gender       *string `json:"gender,omitempty"`
	DateOfBirth  string  `json:"dateOfBirth"`
	MemberSince  *string `json:"memberSince,omitempty"`
	IsStudent    *bool   `json:"isStudent,omitempty"`
	EmailAddress string  `json:"emailAddress"`
	Interests    []byte  `json:"interests,omitempty"`
}

// MemberWriteModel represents a member after a write operation (create/update).
type MemberWriteModel struct {
	ID           int32           `json:"id"`
	Version      int32           `json:"version"`
	Username     string          `json:"username"`
	FirstName    string          `json:"firstName"`
	LastName     string          `json:"lastName"`
	Gender       *string         `json:"gender,omitempty"`
	DateOfBirth  string          `json:"dateOfBirth"`
	MemberSince  *string         `json:"memberSince,omitempty"`
	IsStudent    *bool           `json:"isStudent,omitempty"`
	EmailAddress string          `json:"emailAddress"`
	Interests    json.RawMessage `json:"interests,omitempty"`
	Generated    string          `json:"generated"`
	Updated      string          `json:"updated"`
}

// Create inserts a new member and returns the created record.
func (s *DefaultMemberWriteService) Create(ctx context.Context, cmd CreateMemberCommand) (MemberWriteModel, error) {
	if err := validateCreateMemberCommand(cmd); err != nil {
		return MemberWriteModel{}, err
	}

	params := repository.CreateMemberParams{
		Username:     strings.TrimSpace(cmd.Username),
		FirstName:    strings.TrimSpace(cmd.FirstName),
		LastName:     strings.TrimSpace(cmd.LastName),
		Gender:       mapStringToGender(cmd.Gender),
		DateOfBirth:  parseDate(cmd.DateOfBirth),
		MemberSince:  parseOptionalDate(cmd.MemberSince),
		IsStudent:    mapOptionalBool(cmd.IsStudent),
		EmailAddress: strings.TrimSpace(cmd.EmailAddress),
		Interests:    cmd.Interests,
	}

	member, err := s.members.Create(ctx, params)
	if err != nil {
		return MemberWriteModel{}, err
	}

	return mapMemberWriteModel(member), nil
}

// Update modifies an existing member and returns the updated record.
func (s *DefaultMemberWriteService) Update(ctx context.Context, cmd UpdateMemberCommand) (MemberWriteModel, error) {
	if err := validateUpdateMemberCommand(cmd); err != nil {
		return MemberWriteModel{}, err
	}

	params := repository.UpdateMemberParams{
		ID:           cmd.ID,
		Version:      cmd.Version,
		Username:     strings.TrimSpace(cmd.Username),
		FirstName:    strings.TrimSpace(cmd.FirstName),
		LastName:     strings.TrimSpace(cmd.LastName),
		Gender:       mapStringToGender(cmd.Gender),
		DateOfBirth:  parseDate(cmd.DateOfBirth),
		MemberSince:  parseOptionalDate(cmd.MemberSince),
		IsStudent:    mapOptionalBool(cmd.IsStudent),
		EmailAddress: strings.TrimSpace(cmd.EmailAddress),
		Interests:    cmd.Interests,
	}

	member, err := s.members.Update(ctx, params)
	if errors.Is(err, repository.ErrOptimisticLockStale) {
		return MemberWriteModel{}, ErrOptimisticLockFailed
	}
	if errors.Is(err, repository.ErrMemberNotFound) {
		return MemberWriteModel{}, ErrMemberNotFound
	}
	if err != nil {
		return MemberWriteModel{}, err
	}

	return mapMemberWriteModel(member), nil
}

// Delete removes a member by ID.
func (s *DefaultMemberWriteService) Delete(ctx context.Context, id int32) error {
	if id <= 0 {
		return ErrInvalidMemberInput
	}

	err := s.members.Delete(ctx, id)
	if errors.Is(err, repository.ErrMemberNotFound) {
		return ErrMemberNotFound
	}
	return err
}