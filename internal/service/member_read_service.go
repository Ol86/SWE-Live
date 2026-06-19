package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
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

// ErrInvalidMemberQuery is returned when member query pagination values are invalid.
var ErrInvalidMemberQuery = errors.New("invalid member query")

// MemberReadService provides read-only operations for library members.
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

// normalizeMemberQuery validates pagination and converts empty string filters into nil filters.
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

// normalizeStringFilter trims a string filter and returns nil for empty values.
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

func hasStringFilter(value *string) bool {
	return normalizeStringFilter(value) != nil
}

func isGetAll(filter repository.MemberFilter) bool {
	return filter.Username == nil && filter.EmailAddress == nil && filter.LastName == nil
}

// mapMemberReadModel converts a repository member into the service read model.
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

// mapGender converts a nullable database gender into an optional string.
func mapGender(gender sqlc.NullLibraryGender) *string {
	if !gender.Valid {
		return nil
	}
	value := string(gender.LibraryGender)
	return &value
}

// mapOptionalBool converts a nullable PostgreSQL boolean into an optional bool.
func mapOptionalBool(value pgtype.Bool) *bool {
	if !value.Valid {
		return nil
	}
	return &value.Bool
}

// formatDate formats a required PostgreSQL date as an ISO date string.
func formatDate(value pgtype.Date) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(dateLayout)
}

// formatOptionalDate formats a nullable PostgreSQL date as an optional ISO date string.
func formatOptionalDate(value pgtype.Date) *string {
	if !value.Valid {
		return nil
	}
	formatted := value.Time.Format(dateLayout)
	return &formatted
}

// formatTimestamp formats a PostgreSQL timestamp as an RFC3339 UTC string.
func formatTimestamp(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}
	return value.Time.UTC().Format(time.RFC3339)
}

// mapRawJSON keeps non-empty JSONB values as raw JSON for response encoding.
func mapRawJSON(value []byte) json.RawMessage {
	if len(value) == 0 || string(value) == "null" {
		return nil
	}
	return json.RawMessage(value)
}
