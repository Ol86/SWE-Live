package service

import (
	"encoding/json"
	"strings"
	"time"

	"SWE-Live/internal/db/sqlc"
	"SWE-Live/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

// Mapping functions

func mapMemberWriteModel(member repository.Member) MemberWriteModel {
	return MemberWriteModel{
		ID:           member.ID,
		Version:      member.Version,
		Username:     member.Username,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		Gender:       mapGender(member.Gender),
		DateOfBirth:  formatDate(member.DateOfBirth),
		MemberSince:  formatOptionalDate(member.MemberSince),
		IsStudent:    unmapOptionalBool(member.IsStudent),
		EmailAddress: member.EmailAddress,
		Interests:    mapRawJSON(member.Interests),
		Generated:    formatTimestamp(member.Generated),
		Updated:      formatTimestamp(member.Updated),
	}
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
		IsStudent:    unmapOptionalBool(member.IsStudent),
		EmailAddress: member.EmailAddress,
		Interests:    mapRawJSON(member.Interests),
		Generated:    formatTimestamp(member.Generated),
		Updated:      formatTimestamp(member.Updated),
	}
}

func mapStringToGender(gender *string) sqlc.NullLibraryGender {
	if gender == nil {
		return sqlc.NullLibraryGender{}
	}
	trimmed := stringToUpper(stringTrimSpace(*gender))
	return sqlc.NullLibraryGender{
		LibraryGender: sqlc.LibraryGender(trimmed),
		Valid:         true,
	}
}

func mapOptionalBool(value *bool) pgtype.Bool {
	if value == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *value, Valid: true}
}

func unmapOptionalBool(value pgtype.Bool) *bool {
	if !value.Valid {
		return nil
	}
	return &value.Bool
}

func mapGender(gender sqlc.NullLibraryGender) *string {
	if !gender.Valid {
		return nil
	}
	value := string(gender.LibraryGender)
	return &value
}

func formatDate(value pgtype.Date) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(DateLayout)
}

func formatOptionalDate(value pgtype.Date) *string {
	if !value.Valid {
		return nil
	}
	formatted := value.Time.Format(DateLayout)
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

// Helper functions
func stringTrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func stringToUpper(s string) string {
	return strings.ToUpper(s)
}
