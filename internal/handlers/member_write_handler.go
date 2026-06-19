package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"SWE-Live/internal/service"
)

// MemberWriteHandler handles HTTP requests for member write operations.
type MemberWriteHandler struct {
	service service.MemberWriteService
}

// NewMemberWriteHandler creates a new member write handler.
func NewMemberWriteHandler(svc service.MemberWriteService) *MemberWriteHandler {
	return &MemberWriteHandler{service: svc}
}

// CreateMember handles POST requests to create a new member.
func (h *MemberWriteHandler) CreateMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cmd service.CreateMemberCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	member, err := h.service.Create(r.Context(), cmd)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
}

// UpdateMember handles PUT requests to update an existing member.
func (h *MemberWriteHandler) UpdateMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cmd service.UpdateMemberCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	member, err := h.service.Update(r.Context(), cmd)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// DeleteMember handles DELETE requests to remove a member.
func (h *MemberWriteHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), int32(id)); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleServiceError converts service errors to HTTP responses.
func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidMemberInput):
		fallthrough
	case errors.Is(err, service.ErrInvalidEmail):
		fallthrough
	case errors.Is(err, service.ErrInvalidUsername):
		fallthrough
	case errors.Is(err, service.ErrInvalidDateOfBirth):
		fallthrough
	case errors.Is(err, service.ErrInvalidInterests):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, service.ErrMemberNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, service.ErrOptimisticLockFailed):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
