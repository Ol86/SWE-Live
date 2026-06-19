package repository

import (
	"context"
	"errors"

	"library/internal/db/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	// ErrMemberNotFound is returned when no member exists for the requested lookup.
	ErrMemberNotFound      = errors.New("member not found")
	// ErrOptimisticLockStale is returned when an update targets an outdated member version.
	ErrOptimisticLockStale = errors.New("member was changed concurrently")
)

// Member represents a library member as stored in the database.
type Member = sqlc.LibraryMember

// CreateMemberParams contains the values required to create a member.
type CreateMemberParams = sqlc.CreateMemberParams

// UpdateMemberParams contains the values required to update a member.
type UpdateMemberParams = sqlc.UpdateMemberParams

// MemberFilter contains optional search criteria and pagination settings for member lookups.
type MemberFilter struct {
	// Username filters members by username when set.
	Username     *string
	// EmailAddress filters members by email address when set.
	EmailAddress *string
	// LastName filters members by last name when set.
	LastName     *string
	// Limit caps the number of returned members. Values less than or equal to zero use the default limit.
	Limit        int32
	// Offset skips the given number of matching members before returning results.
	Offset       int32
}

// MemberRepository defines persistence operations for library members.
type MemberRepository interface {
	// FindByID returns the member with the given ID or ErrMemberNotFound if it does not exist.
	FindByID(ctx context.Context, id int32) (Member, error)
	// Find returns members that match the supplied filter.
	Find(ctx context.Context, filter MemberFilter) ([]Member, error)
	// Create inserts a new member and returns the stored record.
	Create(ctx context.Context, member CreateMemberParams) (Member, error)
	// Update modifies an existing member and returns the updated record.
	Update(ctx context.Context, member UpdateMemberParams) (Member, error)
	// Delete removes the member with the given ID or returns ErrMemberNotFound if it does not exist.
	Delete(ctx context.Context, id int32) error
}

// SQLMemberRepository implements MemberRepository using sqlc-generated queries.
type SQLMemberRepository struct {
	queries *sqlc.Queries
}

// NewMemberRepository creates a SQL-backed member repository.
func NewMemberRepository(db sqlc.DBTX) *SQLMemberRepository {
	return &SQLMemberRepository{
		queries: sqlc.New(db),
	}
}

// FindByID returns the member with the given ID.
func (r *SQLMemberRepository) FindByID(ctx context.Context, id int32) (Member, error) {
	member, err := r.queries.GetMemberByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Member{}, ErrMemberNotFound
	}
	return member, err
}

// Find returns members matching the provided filter, using a default limit when none is supplied.
func (r *SQLMemberRepository) Find(ctx context.Context, filter MemberFilter) ([]Member, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	return r.queries.SearchMembers(ctx, sqlc.SearchMembersParams{
		Username:     textParam(filter.Username),
		EmailAddress: textParam(filter.EmailAddress),
		LastName:     textParam(filter.LastName),
		Limit:        limit,
		Offset:       filter.Offset,
	})
}

// Create inserts a member and returns the created database record.
func (r *SQLMemberRepository) Create(ctx context.Context, member CreateMemberParams) (Member, error) {
	return r.queries.CreateMember(ctx, member)
}

// Update modifies a member and returns ErrOptimisticLockStale when the stored version has changed.
func (r *SQLMemberRepository) Update(ctx context.Context, member UpdateMemberParams) (Member, error) {
	updated, err := r.queries.UpdateMember(ctx, member)
	if errors.Is(err, pgx.ErrNoRows) {
		return Member{}, ErrOptimisticLockStale
	}
	return updated, err
}

// Delete removes a member by ID.
func (r *SQLMemberRepository) Delete(ctx context.Context, id int32) error {
	rowsAffected, err := r.queries.DeleteMember(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrMemberNotFound
	}
	return nil
}

// textParam converts an optional string pointer into a nullable PostgreSQL text value.
func textParam(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}
