package service

import (
	"context"
	"encoding/json"

	"SWE-Live/internal/repository"
)

type MemberReadService interface {
	GetByID(ctx context.Context, id int32) (MemberReadModel, error)
	GetByQueryParam(ctx context.Context, query MemberQuery) ([]MemberReadModel, error)
}

type DefaultMemberReadService struct {
	members repository.MemberRepository
}

func NewMemberReadService(members repository.MemberRepository) *DefaultMemberReadService {
	return &DefaultMemberReadService{members: members}
}

type MemberQuery struct {
	Username     *string
	EmailAddress *string
	LastName     *string
	Limit        int32
	Offset       int32
}

type MemberReadModel struct {
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

func (s *DefaultMemberReadService) GetByID(ctx context.Context, id int32) (MemberReadModel, error) {
	member, err := s.members.FindByID(ctx, id)
	if err != nil {
		return MemberReadModel{}, err
	}
	return mapMemberReadModel(member), nil
}

func (s *DefaultMemberReadService) GetByQueryParam(ctx context.Context, query MemberQuery) ([]MemberReadModel, error) {
	filter, err := normalizeMemberQuery(query)
	if err != nil {
		return nil, err
	}

	members, err := s.members.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]MemberReadModel, 0, len(members))
	for _, member := range members {
		result = append(result, mapMemberReadModel(member))
	}
	return result, nil
}
