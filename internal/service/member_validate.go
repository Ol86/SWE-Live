package service

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"SWE-Live/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

// Exported constants for use across service files
const DateLayout = "2006-01-02"
const (
	DefaultMemberReadLimit int32 = 20
	MaxMemberReadLimit     int32 = 100
)

var (
	ErrInvalidMemberInput = errors.New("invalid member input")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidDateOfBirth = errors.New("invalid date of birth")
	ErrInvalidInterests   = errors.New("invalid interests JSON")
	ErrInvalidMemberQuery = errors.New("invalid member query")
)

// Email regex pattern for basic validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Member validation functions

func validateCreateMemberCommand(cmd CreateMemberCommand) error {
	if strings.TrimSpace(cmd.Username) == "" {
		return ErrInvalidUsername
	}
	if strings.TrimSpace(cmd.FirstName) == "" {
		return ErrInvalidMemberInput
	}
	if strings.TrimSpace(cmd.LastName) == "" {
		return ErrInvalidMemberInput
	}
	if !isValidEmail(cmd.EmailAddress) {
		return ErrInvalidEmail
	}
	if !isValidDateFormat(cmd.DateOfBirth) {
		return ErrInvalidDateOfBirth
	}
	if cmd.MemberSince != nil && !isValidDateFormat(*cmd.MemberSince) {
		return ErrInvalidDateOfBirth
	}
	if cmd.Gender != nil && !isValidGender(*cmd.Gender) {
		return ErrInvalidMemberInput
	}
	if len(cmd.Interests) > 0 && !isValidJSON(cmd.Interests) {
		return ErrInvalidInterests
	}
	return nil
}

func validateUpdateMemberCommand(cmd UpdateMemberCommand) error {
	if cmd.ID <= 0 {
		return ErrInvalidMemberInput
	}
	if cmd.Version <= 0 {
		return ErrInvalidMemberInput
	}
	return validateCreateMemberCommand(CreateMemberCommand{
		Username:     cmd.Username,
		FirstName:    cmd.FirstName,
		LastName:     cmd.LastName,
		Gender:       cmd.Gender,
		DateOfBirth:  cmd.DateOfBirth,
		MemberSince:  cmd.MemberSince,
		IsStudent:    cmd.IsStudent,
		EmailAddress: cmd.EmailAddress,
		Interests:    cmd.Interests,
	})
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

func isValidDateFormat(date string) bool {
	_, err := time.Parse(DateLayout, strings.TrimSpace(date))
	return err == nil
}

func isValidGender(gender string) bool {
	g := strings.ToUpper(strings.TrimSpace(gender))
	return g == "MALE" || g == "FEMALE" || g == "DIVERSE"
}

func isValidJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}

// Member query validation functions

func normalizeMemberQuery(query MemberQuery) (repository.MemberFilter, error) {
	if query.Limit < 0 || query.Offset < 0 || query.Limit > MaxMemberReadLimit {
		return repository.MemberFilter{}, ErrInvalidMemberQuery
	}

	limit := query.Limit
	if limit == 0 {
		limit = DefaultMemberReadLimit
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

func hasStringFilter(value *string) bool {
	return normalizeStringFilter(value) != nil
}

func isGetAll(filter repository.MemberFilter) bool {
	return filter.Username == nil && filter.EmailAddress == nil && filter.LastName == nil
}

// Date parsing helpers

func parseDate(dateStr string) pgtype.Date {
	date, err := time.Parse(DateLayout, strings.TrimSpace(dateStr))
	if err != nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: date, Valid: true}
}

func parseOptionalDate(dateStr *string) pgtype.Date {
	if dateStr == nil {
		return pgtype.Date{}
	}
	return parseDate(*dateStr)
}
