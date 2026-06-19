//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"SWE-Live/internal/db/sqlc"
	"SWE-Live/internal/repository"
	"SWE-Live/test/integration/testutil"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestMemberReadOperations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	postgres := testutil.StartPostgres(t, ctx)
	pool := testutil.OpenPool(t, ctx, postgres.DSN)
	testutil.ApplySchema(t, ctx, pool)

	members := repository.NewMemberRepository(pool)

	alice := createMember(t, ctx, members, repository.CreateMemberParams{
		Username:     "alice.reader",
		FirstName:    "Alice",
		LastName:     "Miller",
		Gender:       gender(sqlc.LibraryGenderFEMALE),
		DateOfBirth:  date("1995-03-14"),
		MemberSince:  date("2022-01-10"),
		IsStudent:    boolean(true),
		EmailAddress: "alice@example.com",
		Interests:    rawJSON(t, `["fiction","history"]`),
	})
	bob := createMember(t, ctx, members, repository.CreateMemberParams{
		Username:     "bob.writer",
		FirstName:    "Bob",
		LastName:     "Miller",
		DateOfBirth:  date("1988-08-21"),
		EmailAddress: "bob@sample.net",
	})
	clara := createMember(t, ctx, members, repository.CreateMemberParams{
		Username:     "clara.reader",
		FirstName:    "Clara",
		LastName:     "Adams",
		Gender:       gender(sqlc.LibraryGenderDIVERSE),
		DateOfBirth:  date("2001-11-02"),
		IsStudent:    boolean(false),
		EmailAddress: "clara@example.org",
	})

	t.Run("find by id", func(t *testing.T) {
		got, err := members.FindByID(ctx, alice.ID)
		if err != nil {
			t.Fatalf("find member by id: %v", err)
		}

		assertMember(t, got, alice.ID, "alice.reader", "Alice", "Miller", "alice@example.com")
		if got.Version != 0 {
			t.Fatalf("expected initial version 0, got %d", got.Version)
		}
		if string(got.Interests) != `["fiction", "history"]` {
			t.Fatalf("expected normalized interests JSON, got %s", got.Interests)
		}
	})

	t.Run("missing member by id", func(t *testing.T) {
		_, err := members.FindByID(ctx, 999999)
		if !errors.Is(err, repository.ErrMemberNotFound) {
			t.Fatalf("expected ErrMemberNotFound, got %v", err)
		}
	})

	t.Run("find all defaults to id order", func(t *testing.T) {
		got, err := members.Find(ctx, repository.MemberFilter{})
		if err != nil {
			t.Fatalf("find all members: %v", err)
		}

		assertMemberIDs(t, got, alice.ID, bob.ID, clara.ID)
	})

	t.Run("filter by username fragment", func(t *testing.T) {
		username := "reader"
		got, err := members.Find(ctx, repository.MemberFilter{Username: &username})
		if err != nil {
			t.Fatalf("find members by username: %v", err)
		}

		assertMemberIDs(t, got, alice.ID, clara.ID)
	})

	t.Run("filter by email fragment", func(t *testing.T) {
		email := "sample"
		got, err := members.Find(ctx, repository.MemberFilter{EmailAddress: &email})
		if err != nil {
			t.Fatalf("find members by email: %v", err)
		}

		assertMemberIDs(t, got, bob.ID)
	})

	t.Run("filter by last name is case insensitive", func(t *testing.T) {
		lastName := "miller"
		got, err := members.Find(ctx, repository.MemberFilter{LastName: &lastName})
		if err != nil {
			t.Fatalf("find members by last name: %v", err)
		}

		assertMemberIDs(t, got, alice.ID, bob.ID)
	})

	t.Run("limit and offset", func(t *testing.T) {
		got, err := members.Find(ctx, repository.MemberFilter{Limit: 1, Offset: 1})
		if err != nil {
			t.Fatalf("find members with pagination: %v", err)
		}

		assertMemberIDs(t, got, bob.ID)
	})
}

func createMember(t *testing.T, ctx context.Context, members repository.MemberRepository, params repository.CreateMemberParams) repository.Member {
	t.Helper()

	member, err := members.Create(ctx, params)
	if err != nil {
		t.Fatalf("create member %q: %v", params.Username, err)
	}
	return member
}

func assertMember(t *testing.T, got repository.Member, id int32, username string, firstName string, lastName string, email string) {
	t.Helper()

	if got.ID != id {
		t.Fatalf("expected id %d, got %d", id, got.ID)
	}
	if got.Username != username || got.FirstName != firstName || got.LastName != lastName || got.EmailAddress != email {
		t.Fatalf("unexpected member: %+v", got)
	}
}

func assertMemberIDs(t *testing.T, got []repository.Member, want ...int32) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %d members, got %d: %+v", len(want), len(got), got)
	}

	for i := range want {
		if got[i].ID != want[i] {
			t.Fatalf("member %d: expected id %d, got %d", i, want[i], got[i].ID)
		}
	}
}

func date(value string) pgtype.Date {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		panic(err)
	}
	return pgtype.Date{Time: parsed, Valid: true}
}

func gender(value sqlc.LibraryGender) sqlc.NullLibraryGender {
	return sqlc.NullLibraryGender{LibraryGender: value, Valid: true}
}

func boolean(value bool) pgtype.Bool {
	return pgtype.Bool{Bool: value, Valid: true}
}

func rawJSON(t *testing.T, value string) []byte {
	t.Helper()

	var raw json.RawMessage
	if err := json.Unmarshal([]byte(value), &raw); err != nil {
		t.Fatalf("parse test JSON: %v", err)
	}
	return raw
}
