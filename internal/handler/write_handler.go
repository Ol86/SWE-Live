package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"SWE-Live/internal/service"

	"github.com/gin-gonic/gin"
)

// MemberWriteHandler handles HTTP requests for member write operations.
type MemberWriteHandler struct {
	service service.MemberWriteService
}

// NewMemberWriteHandler creates a new member write handler.
func NewMemberWriteHandler(svc service.MemberWriteService) *MemberWriteHandler {
	return &MemberWriteHandler{service: svc}
}

// RegisterRoutes registers the member write routes on the supplied router group.
func (h *MemberWriteHandler) RegisterRoutes(router gin.IRoutes) {
	slog.Debug("Registering member write routes")
	router.POST("/members", h.CreateMember)
	router.PUT("/members", h.UpdateMember)
	router.DELETE("/members", h.DeleteMember)
}

// CreateMember handles POST /members to create a new member.
func (h *MemberWriteHandler) CreateMember(ctx *gin.Context) {
	slog.DebugContext(ctx.Request.Context(), "Handling create member request",
		"path", ctx.FullPath(),
	)

	var cmd service.CreateMemberCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		slog.DebugContext(ctx.Request.Context(), "Rejected create member request because body is invalid",
			"error", err,
		)
		writeError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	member, err := h.service.Create(ctx.Request.Context(), cmd)
	if err != nil {
		h.writeServiceError(ctx, err, "create member")
		return
	}

	slog.DebugContext(ctx.Request.Context(), "Completed create member request",
		"status", http.StatusCreated,
	)
	ctx.JSON(http.StatusCreated, member)
}

// UpdateMember handles PUT /members to update an existing member.
func (h *MemberWriteHandler) UpdateMember(ctx *gin.Context) {
	slog.DebugContext(ctx.Request.Context(), "Handling update member request",
		"path", ctx.FullPath(),
	)

	var cmd service.UpdateMemberCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		slog.DebugContext(ctx.Request.Context(), "Rejected update member request because body is invalid",
			"error", err,
		)
		writeError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	member, err := h.service.Update(ctx.Request.Context(), cmd)
	if err != nil {
		h.writeServiceError(ctx, err, "update member")
		return
	}

	slog.DebugContext(ctx.Request.Context(), "Completed update member request",
		"status", http.StatusOK,
	)
	ctx.JSON(http.StatusOK, member)
}

// DeleteMember handles DELETE /members and expects an 'id' query parameter.
func (h *MemberWriteHandler) DeleteMember(ctx *gin.Context) {
	slog.DebugContext(ctx.Request.Context(), "Handling delete member request",
		"path", ctx.FullPath(),
		"raw_member_id", ctx.Query("id"),
	)

	id, ok := parseInt32QueryParam(ctx, "id")
	if !ok {
		slog.DebugContext(ctx.Request.Context(), "Rejected delete member request because id is invalid",
			"raw_member_id", ctx.Query("id"),
		)
		writeError(ctx, http.StatusBadRequest, "invalid member id")
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), id); err != nil {
		h.writeServiceError(ctx, err, "delete member", "member_id", id)
		return
	}

	slog.DebugContext(ctx.Request.Context(), "Completed delete member request",
		"member_id", id,
		"status", http.StatusNoContent,
	)
	ctx.Status(http.StatusNoContent)
}

func (h *MemberWriteHandler) writeServiceError(ctx *gin.Context, err error, operation string, attrs ...any) {
	switch {
	case errors.Is(err, service.ErrInvalidMemberInput),
		errors.Is(err, service.ErrInvalidEmail),
		errors.Is(err, service.ErrInvalidUsername),
		errors.Is(err, service.ErrInvalidDateOfBirth),
		errors.Is(err, service.ErrInvalidInterests):
		slog.DebugContext(ctx.Request.Context(), "Member write request returned invalid input",
			append([]any{"operation", operation, "error", err}, attrs...)...,
		)
		writeError(ctx, http.StatusBadRequest, err.Error())

	case errors.Is(err, service.ErrMemberNotFound):
		slog.DebugContext(ctx.Request.Context(), "Member write request returned not found",
			append([]any{"operation", operation, "error", err}, attrs...)...,
		)
		writeError(ctx, http.StatusNotFound, err.Error())

	case errors.Is(err, service.ErrOptimisticLockFailed):
		slog.DebugContext(ctx.Request.Context(), "Member write request returned optimistic lock failure",
			append([]any{"operation", operation, "error", err}, attrs...)...,
		)
		writeError(ctx, http.StatusConflict, err.Error())

	default:
		slog.ErrorContext(ctx.Request.Context(), "Member write request failed",
			append([]any{"operation", operation, "error", err}, attrs...)...,
		)
		writeError(ctx, http.StatusInternalServerError, "internal server error")
	}
}

func parseInt32QueryParam(ctx *gin.Context, name string) (int32, bool) {
	raw := ctx.Query(name)
	if raw == "" {
		return 0, false
	}

	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, false
	}
	return int32(value), true
}
