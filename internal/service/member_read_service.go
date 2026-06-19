package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"SWE-Live/internal/repository"
)

type MemberReadService interface {
	// GetByID returns a single member by its database ID.
	GetByID(ctx context.Context, id int32) (MemberReadModel, error)
	// GetByQueryParam returns members matching the supplied query parameters.
	GetByQueryParam(ctx context.Context, query MemberQuery) ([]MemberReadModel, error)
}

// DefaultMemberReadService implements MemberReadService using a member repository.
type DefaultMemberReadService struct {
	members repository.MemberRepository
}

// NewMemberReadService creates a read service backed by the supplied repository.
func NewMemberReadService(members repository.MemberRepository) *DefaultMemberReadService {
	return &DefaultMemberReadService{members: members}
}

// MemberQuery contains optional filter and pagination values for member reads.
type MemberQuery struct {
	// Username filters members by username. Empty values are treated as unset.
	Username *string
	// EmailAddress filters members by email address. Empty values are treated as unset.
	EmailAddress *string
	// LastName filters members by last name. Empty values are treated as unset.
	LastName *string
	// Limit caps the number of returned members. Zero uses the service default.
	Limit int32
	// Offset skips the given number of matching members.
	Offset int32
}

// MemberReadModel is the JSON-ready representation returned by read operations.
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

// GetByID loads a member from the repository and maps it to the read model.
func (s *DefaultMemberReadService) GetByID(ctx context.Context, id int32) (MemberReadModel, error) {
	slog.DebugContext(ctx, "Loading member by id", "member_id", id)

	member, err := s.members.FindByID(ctx, id)
	if err != nil {
		slog.DebugContext(ctx, "Loading member by id failed", "member_id", id, "error", err)
		return MemberReadModel{}, err
	}

	slog.DebugContext(ctx, "Loaded member by id", "member_id", id, "version", member.Version)
	return mapMemberReadModel(member), nil
}

// GetByQueryParam normalizes query parameters, performs the lookup, and maps the results.
// If all query parameters are empty, the repository is called without filters, which acts as getAll.
func (s *DefaultMemberReadService) GetByQueryParam(ctx context.Context, query MemberQuery) ([]MemberReadModel, error) {
	slog.DebugContext(ctx, "Loading members by query parameters",
		"has_username_filter", hasStringFilter(query.Username),
		"has_email_filter", hasStringFilter(query.EmailAddress),
		"has_last_name_filter", hasStringFilter(query.LastName),
		"limit", query.Limit,
		"offset", query.Offset,
	)

	filter, err := normalizeMemberQuery(query)
	if err != nil {
		slog.DebugContext(ctx, "Member query validation failed",
			"limit", query.Limit,
			"offset", query.Offset,
			"error", err,
		)
		return nil, err
	}

	slog.DebugContext(ctx, "Normalized member query",
		"is_get_all", isGetAll(filter),
		"limit", filter.Limit,
		"offset", filter.Offset,
	)

	members, err := s.members.Find(ctx, filter)
	if err != nil {
		slog.DebugContext(ctx, "Loading members by query parameters failed", "error", err)
		return nil, err
	}

	result := make([]MemberReadModel, 0, len(members))
	for _, member := range members {
		result = append(result, mapMemberReadModel(member))
	}

	slog.DebugContext(ctx, "Loaded members by query parameters", "result_count", len(result))
	return result, nil
}
