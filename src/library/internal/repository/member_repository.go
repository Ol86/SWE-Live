package repository

import (
	"context"
	"errors"

	"library/internal/db/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrMemberNotFound      = errors.New("member not found")
	ErrOptimisticLockStale = errors.New("member was changed concurrently")
)

type Member = sqlc.LibraryMember
type CreateMemberParams = sqlc.CreateMemberParams
type UpdateMemberParams = sqlc.UpdateMemberParams

type MemberFilter struct {
	Username     *string
	EmailAddress *string
	LastName     *string
	Limit        int32
	Offset       int32
}

type MemberRepository interface {
	FindByID(ctx context.Context, id int32) (Member, error)
	Find(ctx context.Context, filter MemberFilter) ([]Member, error)
	Create(ctx context.Context, member CreateMemberParams) (Member, error)
	Update(ctx context.Context, member UpdateMemberParams) (Member, error)
	Delete(ctx context.Context, id int32) error
}

type SQLMemberRepository struct {
	queries *sqlc.Queries
}

func NewMemberRepository(db sqlc.DBTX) *SQLMemberRepository {
	return &SQLMemberRepository{
		queries: sqlc.New(db),
	}
}

func (r *SQLMemberRepository) FindByID(ctx context.Context, id int32) (Member, error) {
	member, err := r.queries.GetMemberByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Member{}, ErrMemberNotFound
	}
	return member, err
}

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

func (r *SQLMemberRepository) Create(ctx context.Context, member CreateMemberParams) (Member, error) {
	return r.queries.CreateMember(ctx, member)
}

func (r *SQLMemberRepository) Update(ctx context.Context, member UpdateMemberParams) (Member, error) {
	updated, err := r.queries.UpdateMember(ctx, member)
	if errors.Is(err, pgx.ErrNoRows) {
		return Member{}, ErrOptimisticLockStale
	}
	return updated, err
}

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

func textParam(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}
