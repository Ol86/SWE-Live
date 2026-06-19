package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"SWE-Live/internal/repository"
	"SWE-Live/internal/service"

	"github.com/gin-gonic/gin"
)

// MemberReadHandler exposes HTTP endpoints for read-only member operations.
type MemberReadHandler struct {
	members service.MemberReadService
}

// NewMemberReadHandler creates a read handler backed by the supplied service.
func NewMemberReadHandler(members service.MemberReadService) *MemberReadHandler {
	return &MemberReadHandler{members: members}
}

// RegisterRoutes registers the member read routes on the supplied router group.
func (h *MemberReadHandler) RegisterRoutes(router gin.IRoutes) {
	slog.Debug("Registering member read routes")
	router.GET("/members", h.GetByQueryParam)
	router.GET("/members/:id", h.GetByID)
}

// GetByID handles GET /members/:id.
func (h *MemberReadHandler) GetByID(ctx *gin.Context) {
	slog.DebugContext(ctx.Request.Context(), "Handling get member by id request",
		"path", ctx.FullPath(),
		"raw_member_id", ctx.Param("id"),
	)

	id, ok := parseInt32PathParam(ctx, "id")
	if !ok {
		slog.DebugContext(ctx.Request.Context(), "Rejected get member by id request because id is invalid",
			"raw_member_id", ctx.Param("id"),
		)
		writeError(ctx, http.StatusBadRequest, "invalid member id")
		return
	}

	member, err := h.members.GetByID(ctx.Request.Context(), id)
	if err != nil {
		h.writeServiceError(ctx, err, "get member by id", "member_id", id)
		return
	}

	slog.DebugContext(ctx.Request.Context(), "Completed get member by id request",
		"member_id", id,
		"status", http.StatusOK,
	)
	ctx.JSON(http.StatusOK, member)
}

// GetByQueryParam handles GET /members and treats an empty query as getAll.
func (h *MemberReadHandler) GetByQueryParam(ctx *gin.Context) {
	slog.DebugContext(ctx.Request.Context(), "Handling get members by query request",
		"path", ctx.FullPath(),
		"has_query", ctx.Request.URL.RawQuery != "",
	)

	query, ok := parseMemberQuery(ctx)
	if !ok {
		slog.DebugContext(ctx.Request.Context(), "Rejected get members by query request because query is invalid",
			"has_query", ctx.Request.URL.RawQuery != "",
		)
		writeError(ctx, http.StatusBadRequest, "invalid member query")
		return
	}

	members, err := h.members.GetByQueryParam(ctx.Request.Context(), query)
	if err != nil {
		h.writeServiceError(ctx, err, "get members by query",
			"limit", query.Limit,
			"offset", query.Offset,
		)
		return
	}

	slog.DebugContext(ctx.Request.Context(), "Completed get members by query request",
		"status", http.StatusOK,
		"result_count", len(members),
	)
	ctx.JSON(http.StatusOK, members)
}

func (h *MemberReadHandler) writeServiceError(ctx *gin.Context, err error, operation string, attrs ...any) {
	switch {
	case errors.Is(err, repository.ErrMemberNotFound):
		slog.DebugContext(ctx.Request.Context(), "Member read request returned not found",
			append([]any{"operation", operation, "error", err}, attrs...)...,
		)
		writeError(ctx, http.StatusNotFound, "member not found")
	case errors.Is(err, service.ErrInvalidMemberQuery):
		slog.DebugContext(ctx.Request.Context(), "Member read request returned invalid query",
			append([]any{"operation", operation, "error", err}, attrs...)...,
		)
		writeError(ctx, http.StatusBadRequest, "invalid member query")
	default:
		slog.ErrorContext(ctx.Request.Context(), "Member read request failed",
			append([]any{"operation", operation, "error", err}, attrs...)...,
		)
		writeError(ctx, http.StatusInternalServerError, "internal server error")
	}
}

func parseMemberQuery(ctx *gin.Context) (service.MemberQuery, bool) {
	limit, ok := parseOptionalInt32Query(ctx, "limit")
	if !ok {
		return service.MemberQuery{}, false
	}

	offset, ok := parseOptionalInt32Query(ctx, "offset")
	if !ok {
		return service.MemberQuery{}, false
	}

	return service.MemberQuery{
		Username:     optionalQuery(ctx, "username"),
		EmailAddress: optionalQueryAlias(ctx, "emailAddress", "email_address"),
		LastName:     optionalQueryAlias(ctx, "lastName", "last_name"),
		Limit:        limit,
		Offset:       offset,
	}, true
}

func parseInt32PathParam(ctx *gin.Context, name string) (int32, bool) {
	value, err := strconv.ParseInt(ctx.Param(name), 10, 32)
	if err != nil {
		return 0, false
	}
	return int32(value), true
}

func parseOptionalInt32Query(ctx *gin.Context, name string) (int32, bool) {
	raw, exists := ctx.GetQuery(name)
	if !exists || raw == "" {
		return 0, true
	}

	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, false
	}
	return int32(value), true
}

func optionalQuery(ctx *gin.Context, name string) *string {
	value, exists := ctx.GetQuery(name)
	if !exists {
		return nil
	}
	return &value
}

func optionalQueryAlias(ctx *gin.Context, preferred string, fallback string) *string {
	value := optionalQuery(ctx, preferred)
	if value != nil {
		return value
	}
	return optionalQuery(ctx, fallback)
}

func writeError(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{"error": message})
}
