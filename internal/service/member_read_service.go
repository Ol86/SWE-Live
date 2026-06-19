package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"SWE-Live/internal/db/sqlc"
	"SWE-Live/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	defaultMemberReadLimit int32 = 20
	maxMemberReadLimit     int32 = 100
	dateLayout                   = "2006-01-02"
)

var ErrInvalidMemberQuery = errors.New("invalid member query")

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

func normalizeMemberQuery(query MemberQuery) (repository.MemberFilter, error) {
	if query.Limit < 0 || query.Offset < 0 || query.Limit > maxMemberReadLimit {
		return repository.MemberFilter{}, ErrInvalidMemberQuery
	}

	limit := query.Limit
	if limit == 0 {
		limit = defaultMemberReadLimit
	}

	return repository.MemberFilter{
		Username:     normalizeStringFilter(query.Username),
		EmailAddress: normalizeStringFilter(query.EmailAddress),
		LastName:     normalizeStringFilter(query.LastName),
		Limit:        limit,
		Offset:       query.Offset,
	}, nil
}

func normalizeStringFilter(value *string) *string {
	if value == nil {
		return nil
	}
	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func mapMemberReadModel(member repository.Member) MemberReadModel {
	return MemberReadModel{
		ID:           member.ID,
		Version:      member.Version,
		Username:     member.Username,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		Gender:       mapGender(member.Gender),
		DateOfBirth:  formatDate(member.DateOfBirth),
		MemberSince:  formatOptionalDate(member.MemberSince),
		IsStudent:    mapOptionalBool(member.IsStudent),
		EmailAddress: member.EmailAddress,
		Interests:    mapRawJSON(member.Interests),
		Generated:    formatTimestamp(member.Generated),
		Updated:      formatTimestamp(member.Updated),
	}
}

func mapGender(gender sqlc.NullLibraryGender) *string {
	if !gender.Valid {
		return nil
	}
	value := string(gender.LibraryGender)
	return &value
}

func mapOptionalBool(value pgtype.Bool) *bool {
	if !value.Valid {
		return nil
	}
	return &value.Bool
}

func formatDate(value pgtype.Date) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(dateLayout)
}

func formatOptionalDate(value pgtype.Date) *string {
	if !value.Valid {
		return nil
	}
	formatted := value.Time.Format(dateLayout)
	return &formatted
}

func formatTimestamp(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}
	return value.Time.UTC().Format(time.RFC3339)
}

func mapRawJSON(value []byte) json.RawMessage {
	if len(value) == 0 || string(value) == "null" {
		return nil
	}
	return json.RawMessage(value)
}
